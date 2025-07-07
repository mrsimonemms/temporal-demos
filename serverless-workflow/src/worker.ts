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

import { NativeConnection, Worker } from '@temporalio/worker';
import { join } from 'node:path';
import * as activities from './activities';
import { ServerlessWorkflow } from './serverless-workflow';

async function run() {
  const wf = await ServerlessWorkflow.load('./workflow.yaml');

  await wf.bundle(join(__dirname, 'workflows.js'));

  const connection = await NativeConnection.connect({
    address: process.env.TEMPORAL_ADDRESS,
  });

  try {
    const worker = await Worker.create({
      connection,
      workflowsPath: require.resolve('./workflows'),
      activities,
      taskQueue: 'serverless-workflow',
    });

    // Start accepting tasks on the queue
    await worker.run();
  } finally {
    // Close the connection once the worker has stopped
    await connection.close();
  }
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
