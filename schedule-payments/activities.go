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
	"context"
	"time"

	"github.com/google/uuid"
)

type activities struct {
	data []PaymentData
}

type SendPaymentResult struct {
	AmountInPence int
	TransactionID uuid.UUID
}

// Simulate a database call that returns all the payments due today
func (a *activities) FindPaymentsForDay(ctx context.Context, startTime, endTime time.Time) ([]PaymentData, error) {
	payments := make([]PaymentData, 0)

	now := time.Now().UTC()

	for _, item := range a.data {
		include := false
		switch item.Schedule {
		case ScheduleWeekly:
			// Check if it's today's weekday
			include = int(now.Weekday()) == item.ScheduleTime
		case ScheduleMonthly:
			// Check if it's today's day of month
			include = now.Day() == item.ScheduleTime
		default:
			// Daily - add to list
			include = true
		}

		if include {
			// Include this payment
			payments = append(payments, item)
		}
	}

	return payments, nil
}

func (a *activities) SendPayment(ctx context.Context, payment PaymentData) (*SendPaymentResult, error) {
	time.Sleep(time.Second * 2)

	return &SendPaymentResult{
		AmountInPence: payment.AmountInPence,
		TransactionID: uuid.New(),
	}, nil
}

func NewActivities(data []PaymentData) (*activities, error) {
	return &activities{
		data: data,
	}, nil
}
