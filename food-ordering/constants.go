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

const OrderFoodTaskQueue = "order-food"

var Queries = struct {
	GET_STATUS string
}{
	GET_STATUS: "GET_STATUS",
}

var Signals = struct {
	CHECKOUT string // Submits order for payment
}{
	CHECKOUT: "CHECKOUT",
}

var Updates = struct {
	ADD_ITEM      string // Adds an item to the order
	REMOVE_ITEM   string // Remove an item from the order
	UPDATE_STATUS string // Restaurant updates status of order
}{
	ADD_ITEM:      "ADD_ITEM",
	REMOVE_ITEM:   "REMOVE_ITEM",
	UPDATE_STATUS: "UPDATE_STATUS",
}
