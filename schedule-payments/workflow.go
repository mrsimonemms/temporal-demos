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

package schedulepayments

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Find payments due today
func FindDuePaymentsWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	year, month, day := time.Now().Date()

	startTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(time.Hour * 24)

	logger.Info("Find payments due", "startTime", startTime, "endTime", endTime)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	var a *activities

	workflow.Sleep(ctx, time.Second*5)

	var payments []PaymentData
	if err := workflow.ExecuteActivity(ctx, a.FindPaymentsForDay, startTime, endTime).Get(ctx, &payments); err != nil {
		logger.Error("Error getting payments due today", "error", err)
		return fmt.Errorf("error getting payments due today: %w", err)
	}

	logger.Debug("Making payments", "count", len(payments))

	workflow.Sleep(ctx, time.Second*5)

	futures := map[workflow.Context]workflow.ChildWorkflowFuture{}

	for i, payment := range payments {
		paymentCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf("%s_payment_%d", workflow.GetInfo(ctx).WorkflowExecution.ID, i),
		})

		// Store as a future so these can be triggered in parallel
		futures[paymentCtx] = workflow.ExecuteChildWorkflow(paymentCtx, MakePayment, payment)
	}

	// Now the workflows are running, wait for the results
	for ctx, workflow := range futures {
		if err := workflow.Get(ctx, nil); err != nil {
			logger.Error("Error making payment", "error", err)
			return fmt.Errorf("error making payment: %w", err)
		}
	}

	return nil
}

func MakePayment(ctx workflow.Context, payment PaymentData) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Make payment")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
	})

	var a *activities
	if err := workflow.ExecuteActivity(ctx, a.SendPayment, payment).Get(ctx, nil); err != nil {
		logger.Error("Error sending payment", "error", err)
		return fmt.Errorf("error getting payments due today: %w", err)
	}

	return nil
}
