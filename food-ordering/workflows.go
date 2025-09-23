/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foodordering

import (
	"fmt"
	"slices"
	"time"

	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, state OrderState) error {
	logger := workflow.GetLogger(ctx)

	// Force to be default state - payment not taken yet
	state.Status = OrderStatusDefault

	var cancel workflow.CancelFunc
	ctx, cancel = workflow.WithCancel(ctx)

	var a *activities

	// Query to return status of basket
	if err := workflow.SetQueryHandler(ctx, Queries.GET_STATUS, func(_ []byte) (OrderState, error) {
		logger.Debug("Returning order status")
		return state, nil
	}); err != nil {
		logger.Error("SetQueryHandler failed.", "error", err, "query", Queries.GET_STATUS)
		return err
	}

	// Add a new item into the basket
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		Updates.ADD_ITEM,
		func(ctx workflow.Context, item OrderProduct) error {
			logger.Info("Adding new item to basket", "productId", item.ProductID, "quantity", item.Quantity)
			state.AddItem(item)

			return nil
		},
		workflow.UpdateHandlerOptions{
			Validator: func(ctx workflow.Context, item OrderProduct) error {
				key := slices.IndexFunc(allProducts, func(i Product) bool {
					return i.ProductID == item.ProductID
				})

				if key == -1 {
					return fmt.Errorf("unknown product id")
				}

				if item.Quantity <= 0 {
					return fmt.Errorf("quantity must be minimum of 1")
				}

				return nil
			},
		},
	); err != nil {
		logger.Error("SetUpdateHandlerWithOptions failed.", "Error", err, "update", Updates.ADD_ITEM)
		return err
	}

	// Remove an item from the basket
	if err := workflow.SetUpdateHandler(
		ctx,
		Updates.REMOVE_ITEM,
		func(ctx workflow.Context, item OrderProduct) error {
			logger.Info("Removing item from the basket", "productId", item.ProductID, "quantity", item.Quantity)
			state.RemoveItem(item)

			return nil
		},
	); err != nil {
		logger.Error("SetUpdateHandlerWithOptions failed.", "Error", err, "update", Updates.REMOVE_ITEM)
		return err
	}

	updateInProgress := false
	// Update the order status - this will come from the restaurant
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		Updates.UPDATE_STATUS,
		func(ctx workflow.Context, input string) error {
			updateInProgress = true
			status, _ := ParseOrderStatus(input)

			logger.Info("Updating order status", "status", status)
			state.Status = status

			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			})

			if state.Status == OrderStatusRejected {
				logger.Info("Order cancelled")
				defer func() {
					fmt.Println("Cancelling")
					cancel()
				}()

				if err := workflow.ExecuteActivity(ctx, a.RefundPayment).Get(ctx, nil); err != nil {
					logger.Error("Error refunding payment", "error", err)
					return fmt.Errorf("error refunding payment: %w", err)
				}
			}

			if err := workflow.ExecuteActivity(ctx, a.SendNotification, state.Status).Get(ctx, nil); err != nil {
				logger.Error("Error notifying of status change", "error", err)
				return fmt.Errorf("error notifying of status change: %w", err)
			}
			updateInProgress = false

			return nil
		},
		workflow.UpdateHandlerOptions{
			Validator: func(ctx workflow.Context, input string) error {
				if _, err := ParseOrderStatus(input); err != nil {
					logger.Debug("Invalid status", "input", input)
					return err
				}

				return nil
			},
		},
	); err != nil {
		logger.Error("SetUpdateHandlerWithOptions failed.", "Error", err, "update", Updates.UPDATE_STATUS)
		return err
	}

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	if err := workflow.ExecuteActivity(ctx, a.TakePayment).Get(ctx, nil); err != nil {
		logger.Error("Error taking payment", "error", err)
		return fmt.Errorf("error taking payment: %w", err)
	}

	// Set order status to pending
	state.Status = OrderStatusPending

	if err := workflow.ExecuteActivity(ctx, a.SendNotification, state.Status).Get(ctx, nil); err != nil {
		logger.Error("Error notifying of status change", "error", err)
		return fmt.Errorf("error notifying of status change: %w", err)
	}

	// Wait for the status to be completed
	if ok, err := workflow.AwaitWithTimeout(ctx, time.Hour, func() bool {
		return state.Status == OrderStatusCompleted && !updateInProgress
	}); err != nil {
		logger.Error("Error waiting for workflow to complete", "error", err)
		return fmt.Errorf("error waiting for workflow to complete: %w", err)
	} else if !ok {
		logger.Error("Await timedout")
		return fmt.Errorf("await timedout")
	}

	return nil
}
