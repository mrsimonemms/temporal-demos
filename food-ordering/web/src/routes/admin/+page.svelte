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
  import { onMount } from 'svelte';

  type Status = {
    status: OrderStatus;
    name: string;
  };

  type nextStatus = {
    next?: Status;
  };

  type previousStatus = {
    previous?: Status;
  };

  type O = IOrder & nextStatus & previousStatus;

  function getPrevNext(status: OrderStatus): nextStatus & previousStatus {
    let next: Status | undefined = undefined;
    let previous: Status | undefined = undefined;

    switch (status) {
      case 'PENDING':
        next = {
          name: 'Accept',
          status: 'ACCEPTED',
        };
        previous = {
          name: 'Reject',
          status: 'REJECTED'
        }
        break;
      case 'ACCEPTED':
        next = {
          name: 'Start cooking',
          status: 'PREPARING',
        };
        break;
      case 'PREPARING':
        next = {
          name: 'Out for delivery',
          status: 'READY',
        };
        break;
      case 'READY':
        next = {
          name: 'Mark delivered',
          status: 'COMPLETED',
        };
        break;
      case 'COMPLETED':
        break;
    }
    return {
      next,
      previous,
    };
  }

  async function getOpenOrders() {
    const response = await fetch(`/api/order`, {
      method: 'GET',
    });

    if (!response.ok) {
      console.log(response);
      err = response.statusText;
      return;
    }

    const o = (await response.json()).orders as IOrder[];

    orders = o.map((item) => ({
      ...item,
      ...getPrevNext(item.state.status),
    }));
  }

  async function updateStatus(orderId: string, status: OrderStatus) {
    const response = await fetch(`/api/order/${orderId}/status`, {
      method: 'POST',
      body: JSON.stringify({
        status,
      }),
    });

    if (!response.ok) {
      console.log(response);
      err = response.statusText;
      return;
    }

    await response.json();

    await getOpenOrders();
  }

  onMount(async () => {
    await getOpenOrders();
    setInterval(async () => {
      await getOpenOrders();
    }, 5000);
  });

  let err: string = $state('');
  let orders: O[] = $state([]);
</script>

Manage your kitchen orders

{#if orders.length === 0}
  <div class="my-5 message">
    <div class="message-body">No open orders</div>
  </div>
{/if}

<div class="columns is-multiline">
  {#each orders as item}
    <div class="column is-4">
      <div class="card">
        <div class="card-header">
          <div class="card-header-title">{item.orderId}</div>
        </div>
        <div class="card-content">
          {#each item.state.products as product}
            <p>{product.name} ({product.quantity})</p>
          {/each}
          <p>
            <strong>Status</strong>:
            <span class="is-capitalized">
              {item.state.status.toLowerCase()}
            </span>
          </p>
        </div>
        <div class="card-footer">
          {#if item.previous && item.previous.status}
            <!-- svelte-ignore a11y_invalid_attribute -->
            <a
              href="#"
              onclick={() => updateStatus(item.orderId, item.previous!.status)}
              class="card-footer-item">{item.previous.name}</a
            >
          {/if}
          {#if item.next && item.next.status}
            <!-- svelte-ignore a11y_invalid_attribute -->
            <a
              href="#"
              onclick={() => updateStatus(item.orderId, item.next!.status)}
              class="card-footer-item">{item.next.name}</a
            >
          {/if}
        </div>
      </div>
    </div>
  {/each}
</div>
