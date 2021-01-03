# TODO

## Features

- ~~finish the Collection in-memory implementation~~
- put a (GraphQL? REST?) API in front of said Collection
- populate Collection with a `sync.Once` scrape at startup
- gin up a TypeScript/React web app FE
  - add *'Watched'* checkboxes for every episode
  - retain check status with [Web Storage](https://caniuse.com/namevalue-storage)

## Structure

- ~~refactor around [Cobra](https://github.com/spf13/cobra)~~
- ~~run config through [Viper](https://github.com/spf13/viper)~~

## Miscellaneous

- `diff --recursive --exclude=.git --unidirectional-new-file . ~/git/github.com/jlucktay/template-go`
