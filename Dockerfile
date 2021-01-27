FROM golang:1.15-alpine as builder
LABEL maintainer="nezerphoenix@gmail.com"

RUN apk add gcc musl-dev

RUN go env -w GOPATH=/go

WORKDIR /go/src/github.com/steve-nzr/goff

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT [ "go", "run", "cmd/login/main.go" ]

EXPOSE 23000/tcp 28000/tcp 5400/tcp
