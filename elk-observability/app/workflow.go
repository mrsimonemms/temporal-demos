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

package elk_observability

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, Hello, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed", "result", result)

	return result, nil
}

func Hello(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	metricsHandler := activity.GetMetricsHandler(ctx)
	metricsHandler = recordHelloStart(metricsHandler, "metrics.Hello")
	var err error
	defer func() {
		recordHelloEnd(metricsHandler, err)
	}()

	logger.Info("Activity", "name", name)

	time.Sleep(time.Second * 5)

	resp := "Hello " + name + "!"
	if name == "error" {
		err = fmt.Errorf("invalid name: %s", name)
		resp = ""
	}

	return resp, err
}
