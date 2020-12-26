FROM golang:1.15 AS build

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/go.jlucktay.dev/arrowverse

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -mod=mod -o /bin/arrowverse

FROM scratch

COPY --from=build /bin/arrowverse /bin/arrowverse

ENTRYPOINT [ "/bin/arrowverse" ]
