# TODO

## Features

- ~~finish the Collection in-memory implementation~~
- put a (GraphQL? REST?) API in front of said Collection
- populate Collection with a `sync.Once` scrape at startup
- gin up a TypeScript/React web app FE
  - add *'Watched'* checkboxes for every episode
  - retain check status with [Web Storage](https://caniuse.com/namevalue-storage)
  - throw [Fiber](https://github.com/gofiber/fiber) into the mix as well?

## Structure

- ~~refactor around [Cobra](https://github.com/spf13/cobra)~~
- ~~run config through [Viper](https://github.com/spf13/viper)~~

## Consistency

- ~~grab episodes from <https://arrowverse.info> as well, and compare~~
  - ~~run comparison with the Google `go-cmp` library~~

## Miscellaneous

- flesh out the README
  - Docker image
- `diff --recursive --exclude=.git --unidirectional-new-file . ~/git/github.com/jlucktay/template-go`

## Logging

- wrap `zap` in [logr](https://github.com/go-logr/logr)

## Publishing

- <https://github.com/actions/starter-workflows/blob/main/ci/docker-image.yml>
- <https://github.com/actions/starter-workflows/blob/main/ci/docker-publish.yml>
  - further reading: <https://docs.github.com/en/free-pro-team@latest/packages/guides/configuring-docker-for-use-with-github-packages>
