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

	"go.temporal.io/sdk/activity"
)

// If you wish to connect any dependencies (eg, database), add in here
type activities struct {
	// Database  *gorm.DB
}

// This is a simple activity that says "hello" to the name given
func (a *activities) SayName(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)

	if name == "" {
		name = "anonymous human"
	}

	return "'Ow am ya " + name + "?", nil
}

// If you need to inject dependencies, pass them in here
// The error response can be useful
func NewActivities( /*dbConnection *gorm.DB*/ ) (*activities, error) {
	return &activities{
		// Database: dbConnection,
	}, nil
}
