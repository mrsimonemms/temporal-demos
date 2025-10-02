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
	"log"

	"github.com/mrsimonemms/temporal-demos/nexus/greeter"
	"github.com/mrsimonemms/temporal-demos/nexus/shared"
	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(envconfig.MustLoadDefaultClientOptions())
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the workflow with the task queue "hackathon"
	w := worker.New(c, greeter.TASK_QUEUE_NAME, worker.Options{})

	// Register Nexus services
	svc := nexus.NewService(shared.ServiceName)
	if err := svc.Register(greeter.HelloWorldOperation); err != nil {
		log.Fatalln("Unable to register Nexus handler", err)
	}
	w.RegisterNexusService(svc)

	// Register the workflows
	w.RegisterWorkflow(greeter.HelloWorldWorkflow)

	// Register the activities - you may need to inject dependencies in here
	activities, err := greeter.NewActivities()
	if err != nil {
		log.Fatalln("Error creating activities", err)
	}
	w.RegisterActivity(activities)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
