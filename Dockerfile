FROM golang:1.15 AS build

# Install common CA certificates
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates=20200601~deb10u1 \
  && apt-get autoremove -y \
  && rm -rf /root/.cache

# Prep for Go build with modules to a static binary
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/go.jlucktay.dev/arrowverse

# This will save Go dependencies in the Docker cache, until/unless they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Do the actual build step itself
RUN GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -o /bin/arrowverse

FROM scratch

# Bring common CA certificates and binary over
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/arrowverse /bin/arrowverse

ENTRYPOINT [ "/bin/arrowverse" ]
