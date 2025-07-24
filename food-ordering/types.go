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

import (
	"fmt"
	"strings"
)

type OrderStatus string

const (
	OrderStatusDefault   OrderStatus = "DEFAULT"   // Order not submitted yet
	OrderStatusPending   OrderStatus = "PENDING"   // Order paid and waiting for restaurant to accept
	OrderStatusAccepted  OrderStatus = "ACCEPTED"  // Restaurant accepted order, but not started work yet
	OrderStatusPreparing OrderStatus = "PREPARING" // Restaurant is cooking your food
	OrderStatusReady     OrderStatus = "READY"     // Food is ready for collection/out for delivery
	OrderStatusCompleted OrderStatus = "COMPLETED" // Food given to a hungry person
)

func ParseOrderStatus(status string) (OrderStatus, error) {
	switch strings.ToUpper(status) {
	case "DEFAULT":
		return OrderStatusDefault, nil
	case "PENDING":
		return OrderStatusPending, nil
	case "ACCEPTED":
		return OrderStatusAccepted, nil
	case "PREPARING":
		return OrderStatusPreparing, nil
	case "READY":
		return OrderStatusReady, nil
	case "COMPLETED":
		return OrderStatusCompleted, nil
	}

	var o OrderStatus
	return o, fmt.Errorf("invalid status: %q", o)
}

type Address struct {
	AddressLine1 string `json:"line1"`
	AddressLine2 string `json:"line2"`
	AddressLine3 string `json:"line3"`
	Town         string `json:"town"`
	County       string `json:"county"`
	PostCode     string `json:"postCode"`
}

type OrderState struct {
	Collection      bool           `json:"collection"`
	DeliveryAddress *Address       `json:"deliveryAddress"`
	Email           string         `json:"email"`
	Products        []OrderProduct `json:"products"`
	Status          OrderStatus    `json:"status"`
}

func (o *OrderState) AddItem(item OrderProduct) {
	// Check if we're updating products
	for i := range o.Products {
		if o.Products[i].ProductID != item.ProductID {
			continue
		}

		o.Products[i].Quantity += item.Quantity
		return
	}

	// Otherwise, add products
	o.Products = append(o.Products, item)
}

func (o *OrderState) RemoveItem(item OrderProduct) {
	for i := range o.Products {
		if o.Products[i].ProductID != item.ProductID {
			continue
		}

		o.Products[i].Quantity -= item.Quantity
		if o.Products[i].Quantity <= 0 {
			o.Products = append(o.Products[:i], o.Products[i+1:]...)
		}
		break
	}
}

type OrderProduct struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type Product struct {
	ProductID int `json:"productId"`
	Name string `json:"name"`
}

func NewOrderState() OrderState {
	return OrderState{
		Products: make([]OrderProduct, 0),
		Status:   OrderStatusDefault,
	}
}
