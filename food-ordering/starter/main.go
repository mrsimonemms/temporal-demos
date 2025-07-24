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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	foodordering "github.com/mrsimonemms/temporal-demos/food-ordering"
	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_ADDRESS"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := "ORDER-" + fmt.Sprintf("%d", time.Now().Unix())

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: foodordering.OrderFoodTaskQueue,
	}

	ctx := context.Background()

	state := foodordering.NewOrderState()
	we, err := c.ExecuteWorkflow(
		ctx,
		workflowOptions,
		foodordering.OrderWorkflow,
		state,
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	time.Sleep(time.Second * 5)

	printState(ctx, c, we)

	handler, err := c.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID:   we.GetID(),
		UpdateName:   foodordering.Updates.ADD_ITEM,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args: []any{
			foodordering.OrderProduct{
				ProductID: 1,
				Quantity:  2,
			},
		},
	})
	if err != nil {
		log.Fatal("Erro issuing ADD_ITEM request")
	}
	if err := handler.Get(ctx, nil); err != nil {
		log.Fatalln("Failed to add to basket", err)
	}

	time.Sleep(time.Second * 5)

	if err := c.SignalWorkflow(
		ctx,
		we.GetID(),
		"",
		foodordering.Signals.CHECKOUT,
		nil,
	); err != nil {
		log.Fatalln("Failed to checkout", err)
	}

	// Wait for end of workflow
	log.Println("Waiting for end of workflow")
	if err := we.Get(ctx, nil); err != nil {
		log.Fatalln("Failed", err)
	}

	printState(ctx, c, we)

	log.Println("Order submitted")
}

func getState(ctx context.Context, c client.Client, we client.WorkflowRun) (*foodordering.OrderState, error) {
	resp, err := c.QueryWorkflow(ctx, we.GetID(), "", foodordering.Queries.GET_STATUS)
	if err != nil {
		log.Fatalln("Failed to query workflow", err)
	}
	var result *foodordering.OrderState
	if err := resp.Get(&result); err != nil {
		return nil, fmt.Errorf("unable to decode state query: %w", err)
	}
	return result, nil
}

func printState(ctx context.Context, c client.Client, we client.WorkflowRun) {
	if state, err := getState(ctx, c, we); err != nil {
		log.Fatalln("Failed to get state", err)
	} else {
		log.Printf("State: %+v\n", *state)
	}
}
