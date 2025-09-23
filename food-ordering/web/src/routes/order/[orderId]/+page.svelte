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
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  async function getOrder() {
    const { orderId } = page.params;

    const response = await fetch(`/api/order/${orderId}/status`, {
      method: 'GET',
    });

    if (!response.ok) {
      console.log(response);
      err = response.statusText;
      return;
    }

    order = (await response.json()) as IOrderState;
  }

  function statusToNumber(status: OrderStatus): number {
    switch (status) {
      case 'PENDING':
        return 1;
      case 'ACCEPTED':
        return 2;
      case 'PREPARING':
        return 3;
      case 'READY':
        return 4;
      case 'COMPLETED':
        return 5;
    }
    return 0;
  }

  onMount(async () => {
    await getOrder();
    setInterval(async () => {
      await getOrder();
    }, 5000);
  });

  let err: string = $state('');
  let order: IOrderState | undefined = $state();
</script>

{#if order}
  {#if order.status === 'COMPLETED'}
    <p class="is-size-2">Enjoy your grub</p>
  {:else if order.status === 'REJECTED'}
    <p class="is-size-2">
      Sorry, we can't do your order - your money has been refunded
    </p>
  {:else}
    <p class="mb-2 is-size-2">
      Order:
      <span class="is-lowercase">{order.status}</span>
    </p>
    <progress
      class="progress is-info is-large"
      value={statusToNumber(order.status)}
      max="5"
    >
    </progress>
  {/if}
{/if}
