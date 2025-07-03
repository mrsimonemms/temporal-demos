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
	"time"

	"github.com/google/uuid"
)

type Schedule string

const (
	ScheduleDaily   = "daily"
	ScheduleWeekly  = "weekly"
	ScheduleMonthly = "monthly"
)

type PaymentData struct {
	// Ignored if daily, weekly is day of week (1=monday), monthly is day of month
	ScheduleTime  int
	Schedule      Schedule
	AmountInPence int
	SenderID      string
	ReceipientID  string
}

// This generates the data. This would ordinarily be a database
// connection
func GenerateData() []PaymentData {
	today := time.Now()
	tomorrow := today.Add(time.Hour * 24)
	yesterday := today.Add(-(time.Hour * 24))

	data := []PaymentData{
		{
			// Daily - due today
			Schedule:      ScheduleDaily,
			AmountInPence: 10000,
		},
		{
			// Weekly - due yesterday
			Schedule:      ScheduleWeekly,
			ScheduleTime:  int(yesterday.Weekday()),
			AmountInPence: 10100,
		},
		{
			// Weekly - due today
			Schedule:      ScheduleWeekly,
			ScheduleTime:  int(today.Weekday()),
			AmountInPence: 10200,
		},
		{
			// Weekly - due tomorrow
			Schedule:      ScheduleWeekly,
			ScheduleTime:  int(tomorrow.Weekday()),
			AmountInPence: 10300,
		},
		{
			// Monthly - due yesterday
			Schedule:      ScheduleWeekly,
			ScheduleTime:  yesterday.Day(),
			AmountInPence: 10400,
		},
		{
			// Monthly - due today
			Schedule:      ScheduleWeekly,
			ScheduleTime:  today.Day(),
			AmountInPence: 10000,
		},
		{
			// Monthly - due tomorrow
			Schedule:      ScheduleWeekly,
			ScheduleTime:  tomorrow.Day(),
			AmountInPence: 10000,
		},
	}

	// Generate IDs of the sender and receiver - this would be customer reference IDs
	for i := range data {
		data[i].ReceipientID = uuid.NewString()
		data[i].SenderID = uuid.NewString()
	}

	return data
}
