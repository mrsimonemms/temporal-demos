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

export const products: IProduct[] = [
  {
    id: 1,
    name: 'Chips',
    price: 3.5,
  },
  {
    id: 2,
    name: 'Battered cod',
    price: 8.75,
  },
  {
    id: 3,
    name: 'Battered haddock',
    price: 9.75,
  },
  {
    id: 4,
    name: 'Curry sauce',
    price: 1.45,
  },
  {
    id: 5,
    name: 'Gravy',
    price: 1.45,
  },
];

export function getProduct(t: number): IProduct | undefined {
  return products.find(({ id }) => id === t);
}
