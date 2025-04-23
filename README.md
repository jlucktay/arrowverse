# Arrowverse episode watch order

[![Coverage Status](https://coveralls.io/repos/github/jlucktay/arrowverse/badge.svg?branch=main)](https://coveralls.io/github/jlucktay/arrowverse?branch=main)
[![godoc](https://img.shields.io/badge/pkg.go.dev-godoc-00ADD8?logo=go)](https://pkg.go.dev/go.jlucktay.dev/arrowverse)

## Subcommand layout

- `arrowverse`: prints usage
  - `scrape`: scrapes data to populate a collection, prints data, exits
  - `api`: scrapes as above, launches an API only, no front end
  - `serve` the whole enchilada (all above, plus a front end)

## Design

Internal layout of the CLI inspired in part by:

- [`gh`](https://github.com/cli/cli)
- [`kubectl`](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubectl)
- [Command Line Interface Guidelines](https://clig.dev)
