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

import "go.temporal.io/sdk/client"

const (
	activityLatency      = "activity_latency"
	activityStartedCount = "activity_started"
	activityFailedCount  = "activity_failed"
	activitySuccessCount = "activity_succeeded"
	activityRunningCount = "activity_running"
)

// recordHelloStart emits metrics at the start of an activity function
func recordHelloStart(
	handler client.MetricsHandler,
	activityType string,
) client.MetricsHandler {
	handler = handler.WithTags(map[string]string{"operation": activityType})
	// Increment the number of started activities
	handler.Counter(activityStartedCount).Inc(1)
	return handler
}

// recordHelloEnd emits metrics at the end of an activity function
func recordHelloEnd(handler client.MetricsHandler, err error) {
	if err != nil {
		// Add to failed count
		handler.Counter(activityFailedCount).Inc(1)
		return
	}
	// Add to success count
	handler.Counter(activitySuccessCount).Inc(1)
}
