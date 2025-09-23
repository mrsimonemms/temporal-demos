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

import { json, type RequestHandler } from '@sveltejs/kit';
import Long from 'long';
import { nanoid } from 'nanoid';

function timestampToDate(ts: {
  seconds?: Long | null;
  nanos?: number | null;
}): Date {
  // Convert Long to number (safe for UNIX seconds < ~285,000 years 😅)
  const seconds = ts.seconds?.toNumber() ?? 0;
  const nanos = ts.nanos ?? 0;

  return new Date(seconds * 1000 + nanos / 1e6);
}

import { getProduct } from '$lib/products';
import { ensureConnection } from '$lib/server/temporal';

export const GET: RequestHandler = async () => {
  const temporal = await ensureConnection();

  // This isn't a great way, but I want to avoid a DB for this demo
  const { executions } = await temporal.workflowService.listWorkflowExecutions({
    namespace: 'default',
    query: 'ExecutionStatus="Running"',
  });

  const orders: IOrder[] = [];

  for (const exec of executions) {
    const orderId = exec.execution?.workflowId ?? 'unknown';

    const handler = temporal.workflow.getHandle(orderId);

    const state = (await handler.query('GET_STATUS')) as IProduct2;

    orders.push({
      orderId,
      state: {
        collection: true,
        products: state.products.map((item) => {
          const p = getProduct(item.productId);
          return {
            id: item.productId,
            name: p?.name ?? '',
            price: p?.price ?? 0,
            quantity: item.quantity,
          };
        }),
        status: state.status as OrderStatus,
      },
      created: exec.startTime ? timestampToDate(exec.startTime) : new Date(),
    });
  }

  return json({
    orders,
  });
};

export const POST: RequestHandler = async ({ request }) => {
  const temporal = await ensureConnection();

  const data = await request.json();
  const workflowId = `order-${nanoid()}`;

  await temporal.workflow.start('OrderWorkflow', {
    taskQueue: 'order-food',
    args: [data],
    workflowId,
  });

  return json({
    orderId: workflowId,
  });
};
