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
	"log"
	"os"
	"time"

	schedulepayments "github.com/mrsimonemms/temporal-demos/schedule-payments"
	"go.temporal.io/sdk/client"
)

const scheduleID = "trigger_payments"

// This script upserts a schedule into Temporal that is designed to run indefintely.
// This might be created by a CI/CD action, a Kubernetes Job or any other method
// of running a script to completion.
//
// This means this schedule will remain in your Temporal instance for it's life
// and deletion is out of the scope of this demo. You MUST manually delete it if
// you are using a long-running Temporal service (eg, Temporal Cloud).
func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_ADDRESS"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	ctx := context.Background()

	log.Printf("Looking for existing schedule: %s", scheduleID)
	schedules, err := c.ScheduleClient().List(ctx, client.ScheduleListOptions{})
	if err != nil {
		log.Fatalln("Unable to list schedule", err)
	}

	for schedules.HasNext() {
		s, err := schedules.Next()
		if err != nil {
			log.Fatalln("Unable to get schedule", err)
		}

		// Find and destroy the schedule
		if s.ID == scheduleID {
			log.Printf("Schedule already exists - deleting it")
			handler := c.ScheduleClient().GetHandle(ctx, scheduleID)

			if err := handler.Delete(ctx); err != nil {
				log.Fatalln("Error deleting schedule", err)
			}
		}
	}

	_, err = c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			// Run every day at 2am - this is what this demo would normally run
			Calendars: []client.ScheduleCalendarSpec{
				{
					Hour: []client.ScheduleRange{
						{
							Start: 2,
						},
					},
				},
			},
			// Run every minute - this demonstrates what the schedule does
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: time.Minute,
				},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			Workflow:  schedulepayments.FindDuePaymentsWorkflow,
			TaskQueue: "payments",
		},
	})
	if err != nil {
		log.Fatalln("Error creating schedule", err)
	}

	log.Println("Schedule configured - goodbye")
}
