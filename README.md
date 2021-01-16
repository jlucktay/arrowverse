# Arrowverse episode watch order

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
