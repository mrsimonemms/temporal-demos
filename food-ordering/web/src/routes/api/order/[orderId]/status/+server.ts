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

import { ensureConnection } from '$lib/server/temporal';
import { json, type RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ params }) => {
  const temporal = await ensureConnection();

  const handler = temporal.workflow.getHandle(params.orderId ?? '');

  // Validate the workflow exists
  await handler.describe();

  return json(await handler.query('GET_STATUS'));
};

export const POST: RequestHandler = async ({ params, request }) => {
  const temporal = await ensureConnection();

  const data = await request.json();
  const handler = temporal.workflow.getHandle(params.orderId ?? '');

  // Validate the workflow exists
  await handler.describe();

  await handler.executeUpdate('UPDATE_STATUS', {
    args: [data.status],
  });

  return json({
    hello: 'world',
  });
};
