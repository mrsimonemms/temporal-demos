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

interface IProduct {
  id: number;
  name: string;
  price: number;
  quantity?: number;
}

interface IProduct2 {
  collection: boolean;
  products: { productId: number; quantity: number }[];
  status: OrderStatus;
}

type OrderStatus =
  | 'DEFAULT' // Order not submitted yet
  | 'PENDING' // Order paid and waiting for restaurant to accept
  | 'ACCEPTED' // Restaurant accepted order, but not started work yet
  | 'PREPARING' // Restaurant is cooking your food
  | 'READY' // Food is ready for collection/out for delivery
  | 'COMPLETED'; // Food given to a hungry person

interface IOrderState {
  collection: boolean;
  products: IProduct[];
  status: OrderStatus;
}

interface IOrder {
  orderId: string;
  state: IOrderState;
  created: Date;
}
