# Nexus

This demonstrates how to set up a Nexus application

<!-- toc -->

* [Overview](#overview)
* [Running](#running)
  * [Greeter](#greeter)
  * [Champions](#champions)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Overview

This is two independent applications:

* Greeter: greets whomever is they're told to, in the finest
  [Black Country dialect](https://en.wikipedia.org/wiki/Black_Country_dialect)
* Champions: finds a random F1 champions name and uses the Greeter to greet them
  (yes, all F1 drivers should be greeted like Nigel Mansell speaks)

Whilst a contrived example, it shows how two independent teams can maintain their
own Temporal workflows, but allow them to interact with one another. The Greeter
app publishes an interface and, so long as that's adhered to, the Champions app
can use it to send their messages.

In this example, they're only responding to the caller. However, this could easily
be extended to allow triggering of a text message, an email or even the sending
of a radio message from Bono to Lewis ("it's Hammertime!").

## Running

> This is only checked against Temporal Cloud, although this can easily be used
> against the Temporal CLI dev server. Create an `.envrc` file inside each app's
> root and configure the following information:
>
> ```sh
> export TEMPORAL_ADDRESS="<address of temporal namespace endpoint>"
> export TEMPORAL_API_KEY="<api key>"
> export TEMPORAL_NAMESPACE="<namespace>"
> export TEMPORAL_USE_TLS=true
> ```

### Greeter

Run the workflow, connecting to your first Temporal namespace

```sh
cd /path/to/nexus/greeter
air
```

### Champions

Run the workflow, connecting to a different Temporal namespace

```sh
cd /path/to/nexus/champions
air
```

Trigger the workflow

```sh
cd /path/to/nexus/champions/starter
go run .
```

After a few seconds, you should see which champion you have greeted.
