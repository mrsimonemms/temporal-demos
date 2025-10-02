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

package greeter

import (
	"context"
	"time"

	"github.com/mrsimonemms/temporal-demos/nexus/shared"
	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/workflow"
)

func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var a *activities

	var result string
	err := workflow.ExecuteActivity(ctx, a.SayName, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

var HelloWorldOperation = temporalnexus.NewWorkflowRunOperation(shared.HelloOperationName, HelloWorldWorkflow, func(ctx context.Context, input string, options nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
	return client.StartWorkflowOptions{
		// Workflow IDs should typically be business meaningful IDs and are used to dedupe workflow starts.
		// For this example, we're using the request ID allocated by Temporal when the caller workflow schedules
		// the operation, this ID is guaranteed to be stable across retries of this operation.
		ID: options.RequestID,
		// Task queue defaults to the task queue this operation is handled on.
	}, nil
})
