# Inspiration:
# - https://devhints.io/makefile
# - https://tech.davis-hansson.com/p/make/

SHELL := bash

# Default - top level rule is what gets run when you run just `make` without specifying a goal/target.
.DEFAULT_GOAL := build

.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --warn-undefined-variables

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later.)
endif
.RECIPEPREFIX = >

image_repository := "jlucktay/arrowverse"

test: tmp/.tests-passed.sentinel
lint: tmp/.linted.sentinel
build: out/image-id
.PHONY: test lint build

bench: tmp/.benchmarks-ran.sentinel
.PHONY: bench

# Clean up the output directories; all the sentinel files go under `tmp`, so this will cause everything to get rebuilt.
clean:
> rm -rf ./tmp
> rm -rf ./out
.PHONY: clean

# Clean up any built Docker images.
clean-docker:
> docker images \
  --filter=reference=$(image_repository) \
  --no-trunc --quiet | sort -f | uniq | xargs -n 1 docker rmi --force
> rm -f ./out/image-id
.PHONY: clean-docker

# Tests - re-run if any Go files have changes since `tmp/.tests-passed.sentinel` was last touched.
tmp/.tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test ./...
> touch $@

# Lint - re-run if the tests have been re-run (and so, by proxy, whenever the source files have changed).
tmp/.linted.sentinel: Dockerfile tmp/.tests-passed.sentinel
> mkdir -p $(@D)
> hadolint Dockerfile
> find . -type f -iname "*.go" -exec gofmt -s -w "{}" +
> go vet ./...
> golangci-lint run
> touch $@

# Docker image - re-build if the lint output is re-run.
out/image-id: Dockerfile tmp/.linted.sentinel
> mkdir -p $(@D)
> image_id="$(image_repository):$$(uuidgen)"
> docker build --tag="$${image_id}" .
> echo "$${image_id}" > out/image-id

# Benchmarks - run enough iterations of each benchmark to take 10 seconds
tmp/.benchmarks-ran.sentinel: tmp/.tests-passed.sentinel
> mkdir -p $(@D)
> go test ./... -bench=. -benchmem -benchtime=10s -run=DoNotRunTests
> touch $@
