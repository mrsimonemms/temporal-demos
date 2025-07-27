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
import type { ServerInit } from '@sveltejs/kit';
import process from 'node:process';

export const init: ServerInit = async () => {
  console.log('Connecting to Temporal service');
  const temporal = await ensureConnection();

  process.on('sveltekit:shutdown', async () => {
    console.log('Closing Temporal connection');
    await temporal.connection.close();
  });
};
