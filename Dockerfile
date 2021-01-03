FROM golang:1.15 AS builder

# Set some shell options for using pipes and such
SHELL [ "/bin/bash", "-euo", "pipefail", "-c" ]

# Install common CA certificates to blag later
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates=20200601~deb10u1 \
  && apt-get autoremove -y \
  && rm -rf /root/.cache

# Don't call any C code (the 'scratch' base image used later won't have any libraries to reference)
ENV CGO_ENABLED=0

# Use Go modules
ENV GO111MODULE=on

# Precompile the entire Go standard library into a Docker cache layer: useful for other projects too!
# cf. https://www.reddit.com/r/golang/comments/hj4n44/improved_docker_go_module_dependency_cache_for/
RUN go install -v std

WORKDIR /go/src/go.jlucktay.dev/arrowverse

# This will save Go dependencies in the Docker cache, until/unless they change
COPY go.mod go.sum ./

# May or may not need this special handling to deal with the protocol buffer dependency
# cf. https://github.com/golang/protobuf/issues/751
# RUN if pb=$(go mod graph | awk '{if ($1 !~ "@") print $2}' | grep "google.golang.org/protobuf"); then \
#   go get -v "${pb/@/\/...@}"; fi

# Download and precompile all third party libraries (protobuf will be dealt with indirectly)
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | grep -v "google.golang.org/protobuf" | xargs go get -v

# Add the sources
COPY . .

# Compile! Should only compile our project since everything else has been precompiled by now, and future
# (re)compilations will leverage the same cached layer(s)
RUN go build -v -o /bin/arrowverse

FROM scratch

# Bring common CA certificates and binary over
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/arrowverse /bin/arrowverse

ENTRYPOINT [ "/bin/arrowverse" ]
