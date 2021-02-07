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
golangci_lint_version := v1.35.2

all: test-cover lint build
test: tmp/.short-tests-passed.sentinel
test-all: tmp/.all-tests-passed.sentinel
test-consistency: tmp/.consistency-tests-passed.sentinel
test-cover: tmp/.cover-tests-passed.sentinel
lint: tmp/.linted.sentinel
build: out/image-id
.PHONY: all test test-all test-consistency test-cover lint build

bench: tmp/.benchmarks-ran.sentinel
.PHONY: bench

# Clean up binaries, test coverage, and the output directories.
# All the sentinel files go under `tmp`, so this will cause everything to get rebuilt.
clean:
> rm -f ./arrowverse ./cover.out
> rm -rf ./tmp ./out
.PHONY: clean

# Clean up any built Docker images.
clean-docker:
> docker images \
  --filter=reference=$(image_repository) \
  --no-trunc --quiet | sort -f | uniq | xargs -n 1 docker rmi --force
> rm -f ./out/image-id
.PHONY: clean-docker

# Clean up any binaries under `hack`.
clean-hack:
> rm -rf ./hack/bin
.PHONY: clean-hack

# Clean all of the things.
clean-all: clean clean-docker clean-hack
.PHONY: clean-all

# Tests - re-run short/all tests if any Go files have changed since the relevant sentinel file was last touched.
tmp/.short-tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -short ./...
> touch $@

tmp/.all-tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -count=1 -race ./...
> touch $@

tmp/.consistency-tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -tags=consistency ./...
> touch $@

tmp/.cover-tests-passed.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test -count=1 -covermode=atomic -coverprofile=cover.out -race ./...
> touch $@

# Lint - re-run if the tests have been re-run (and so, by proxy, whenever the source files have changed).
tmp/.linted.sentinel: Dockerfile .golangci.yaml hack/bin/golangci-lint tmp/.short-tests-passed.sentinel
> mkdir -p $(@D)
> docker run --interactive --rm hadolint/hadolint < Dockerfile
> find . -type f -iname "*.go" -exec gofmt -s -w "{}" +
> go vet ./...
> hack/bin/golangci-lint run
> touch $@

hack/bin/golangci-lint:
> curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
> | sh -s -- -b $(shell pwd)/hack/bin $(golangci_lint_version)

# Docker image - re-build if the lint output is re-run.
out/image-id: Dockerfile tmp/.linted.sentinel
> mkdir -p $(@D)
> image_id="$(image_repository):$(shell uuidgen)"
> DOCKER_BUILDKIT=1 docker build --tag="$${image_id}" .
> echo "$${image_id}" > out/image-id

# Benchmarks - run enough iterations of each benchmark to take a few seconds (default is 1s)
tmp/.benchmarks-ran.sentinel: $(shell find . -type f -iname "*.go")
> mkdir -p $(@D)
> go test ./... -bench=. -benchmem -benchtime=3s -run=DoNotRunTests
> touch $@
