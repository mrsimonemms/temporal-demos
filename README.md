# temporal-demos

Collection of Temporal demos

<!-- toc -->

* [Running samples](#running-samples)
* [Contributing](#contributing)
  * [Open in a container](#open-in-a-container)
  * [Commit style](#commit-style)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Running samples

> It's recommended to [open in a container](https://code.visualstudio.com/docs/devcontainers/containers)
> as all dependencies are installed.

Look in each sample for running instructions. The Temporal service can be run
with `temporal server start-dev`.

If using Dev Containers, the Temporal service will run automatically and the
address is configured to the `TEMPORAL_ADDRESS` environment variable.

## Contributing

### Open in a container

* [Open in a container](https://code.visualstudio.com/docs/devcontainers/containers)

### Commit style

All commits must be done in the [Conventional Commit](https://www.conventionalcommits.org)
format.

```git
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```
