# Schedule Payments

Look for payments due today and schedule them

<!-- toc -->

* [Overview](#overview)
* [Steps to run](#steps-to-run)
  * [Run the worker](#run-the-worker)
  * [Create the schedule](#create-the-schedule)
  * [Trigger a run](#trigger-a-run)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Overview

This is an example of a company that needs to trigger a series of regular
payments. Some payments are daily, some are weekly and others are monthly.
A Temporal schedule is created which runs an activity to find the payments that
are due today and then triggers a child workflow to make each of these payments.

The logic to get due payments is very much a demo - in a production system, you
would want to ensure that you only pull transactions that need to be made from your
database. There is no database in this example, but a simple function to filter
payments not required "today". The purpose is demonstrate how Temporal can
create robust schedules, not to show how to pull things from a database.

## Steps to run

### Run the worker

```sh
go run ./worker
```

The worker is where the workflow is defined.

### Create the schedule

```sh
go run ./schedule
```

Create the schedule.

By default, the schedule runs daily at 2am. From a demo point of view, you don't
want to wait all day for this so it's also triggers every 60 seconds.

### Trigger a run

```sh
go run ./starter
```

This enables you to trigger an individual run, testing out the workflow.
