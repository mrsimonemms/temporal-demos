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
	"time"

	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, state OrderState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, QueryTypes.GET_STATUS, func(input []byte) (OrderState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	checkoutChannel := workflow.GetSignalChannel(ctx, SignalChannels.CHECKOUT_CHANNEL)

	var a *activities

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(checkoutChannel, func(c workflow.ReceiveChannel, _ bool) {
			logger.Info("Checkout triggered")

			// Change state
			state.Status = OrderStatusPending

			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			})

			if err := workflow.ExecuteActivity(ctx, a.DoSomething).Get(ctx, nil); err != nil {
				logger.Error("Error checking out", "error", err)
				return
			}
		})

		selector.Select(ctx)

		if state.Status == OrderStatusPending {
			break
		}
	}

	return nil
}
