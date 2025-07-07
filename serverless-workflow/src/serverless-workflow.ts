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

import { Classes, Specification } from '@serverlessworkflow/sdk';
import fs from 'fs/promises';
import Handlebars from 'handlebars';
import { join } from 'node:path';

export interface IAbstractExpression {
  interpret(): string;
}

export class Task implements IAbstractExpression {
  interpret(): string {
    return '';
  }
}

export interface ITask {
  name: string;
  task: Specification.Task;
}

export class ServerlessWorkflow {
  constructor(private workflow: Specification.Workflow) {
    if (this.workflow.do.length === 0) {
      throw new Error('Task list cannot be empty');
    }
  }

  async bundle(filepath: string): Promise<void> {
    const source = await fs.readFile(join(__dirname, './template.hbs'), {
      encoding: 'utf8',
    });
    const template = Handlebars.compile(source);

    const tasks: ITask[] = [];
    this.workflow.do.forEach((item) => {
      Object.entries(item).forEach(([name, task]) => {
        tasks.push({
          name,
          task,
        });
      });
    });

    console.log(
      template({
        workflow: this.workflow,
        tasks,
      }),
    );
    process.exit(1);
    // console.log(this.generateTemporalWorkflow());
    console.log(filepath);

    return fs.writeFile(filepath, 'hello2', { encoding: 'utf8' });
  }

  static async load(filepath: string): Promise<ServerlessWorkflow> {
    const schemaText = await fs.readFile(filepath, {
      encoding: 'utf8',
    });

    const workflow = Classes.Workflow.deserialize(schemaText);
    workflow.validate();

    return new ServerlessWorkflow(workflow);
  }
}
