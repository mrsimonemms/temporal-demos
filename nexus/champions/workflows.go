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

package champions

import (
	"fmt"
	"time"

	"github.com/mrsimonemms/temporal-demos/nexus/shared"
	"go.temporal.io/sdk/workflow"
)

const (
	endpointName = "simon-nexus-testing"
)

func SayHelloToChampionWorkflow(ctx workflow.Context) (string, error) {
	c := workflow.NewNexusClient(endpointName, shared.ServiceName)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("SayHelloToChampionWorkflow workflow started")

	var a *activities

	// Generate the champion to use locally
	var champion string
	err := workflow.ExecuteActivity(ctx, a.FindChampion).Get(ctx, &champion)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("The champion we're saying hello to:", "champion", champion)

	future := c.ExecuteOperation(ctx, shared.HelloOperationName, champion, workflow.NexusOperationOptions{})

	var result string
	if err := future.Get(ctx, &result); err != nil {
		return "", fmt.Errorf("nexus operation failed: %w", err)
	}

	logger.Info("SayHelloToChampionWorkflow workflow completed.", "name", result)

	return result, nil
}
