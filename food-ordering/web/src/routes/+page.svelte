<!--
  ~ Copyright 2025 Simon Emms <simon@simonemms.com>
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { getProduct, products } from '$lib/products';

  interface IOrder {
    productId: number;
    quantity: number;
  }

  export const order: Record<number, number> = $state({});

  export function addItem(id: number) {
    if (!order[id]) {
      order[id] = 0;
    }

    order[id]++;
  }

  export function removeItem(id: number) {
    if (!order[id]) {
      return;
    }

    order[id]--;

    if (order[id] < 0) {
      order[id] = 0;
    }
  }

  export async function submit() {
    loading = true;
    const response = await fetch('/api/order', {
      method: 'POST',
      body: JSON.stringify({
        collection: true,
        deliveryAddress: null,
        email: 'test@test.com',
        status: 'DEFAULT',
        products: Object.entries(order).map(
          (item): IOrder => ({
            productId: Number(item[0]),
            quantity: item[1],
          }),
        ),
      }),
    });

    if (!response.ok) {
      console.log(response);
      loading = false;
      err = response.statusText;
      return;
    }

    const res = await response.json();

    goto(`/order/${res.orderId}`);
  }

  let loading: boolean = $state(false);
  let err: string = $state('');

  export function getTotal(): Number {
    let total = 0;

    Object.entries(order).forEach((item) => {
      const i = getProduct(Number(item[0]));
      total += (i?.price ?? 0) * item[1];
    });

    return total;
  }
</script>

<div class="columns">
  <div class="column is-half">
    <section class="section">
      <div class="columns is-multiline">
        {#each products as item}
          <div class="column is-half">
            <div class="card">
              <div class="card-header">
                <div class="card-header-title">{item.name}</div>
              </div>
              <div class="card-content">&pound;{item.price.toFixed(2)}</div>
              <div class="card-footer">
                <button
                  onclick={() => removeItem(item.id)}
                  class="card-footer-item"
                  >&minus;
                </button>
                <button
                  onclick={() => addItem(item.id)}
                  class="card-footer-item"
                  >&plus;
                </button>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </section>
  </div>
  <div class="column is-half">
    <p class="title">Order</p>
    <table class="table is-fullwidth">
      <thead>
        <tr>
          <th>Item</th>
          <th>Quantity</th>
          <th>Price</th>
        </tr>
      </thead>
      <tfoot>
        <tr>
          <td></td>
          <td></td>
          <td>&pound;{getTotal().toFixed(2)}</td>
        </tr>
      </tfoot>
      <tbody>
        {#each Object.entries(order) as item}
          {@const i = getProduct(Number(item[0]))}
          <tr>
            <td>{i?.name}</td>
            <td>{item[1]}</td>
            <td>
              &pound;{((i?.price ?? 0) * item[1]).toFixed(2)}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>

    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div onclick={() => submit()} class="button is-primary is-fullwidth">
      Order
    </div>
  </div>
</div>
