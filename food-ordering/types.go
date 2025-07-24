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

package foodordering

type OrderStatus string

const (
	OrderStatusDefault   OrderStatus = "DEFAULT"   // Order not submitted yet
	OrderStatusPending   OrderStatus = "PENDING"   // Order paid and waiting for restaurant to accept
	OrderStatusAccepted  OrderStatus = "ACCEPTED"  // Restaurant accepted order, but not started work yet
	OrderStatusPreparing OrderStatus = "PREPARING" // Restaurant is cooking your food
	OrderStatusReady     OrderStatus = "READY"     // Food is ready for collection/out for delivery
	OrderStatusCompleted OrderStatus = "COMPLETED" // Food given to a hungry person
)

type Address struct {
	AddressLine1 string `json:"line1"`
	AddressLine2 string `json:"line2"`
	AddressLine3 string `json:"line3"`
	Town         string `json:"town"`
	County       string `json:"county"`
	PostCode     string `json:"postCode"`
}

type OrderState struct {
	Collection      bool        `json:"collection"`
	DeliveryAddress *Address    `json:"deliveryAddress"`
	Email           string      `json:"email"`
	Products        []Product   `json:"products"`
	Status          OrderStatus `json:"status"`
}

type Product struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

func NewOrderState() OrderState {
	return OrderState{
		Products: make([]Product, 0),
		Status:   OrderStatusDefault,
	}
}
