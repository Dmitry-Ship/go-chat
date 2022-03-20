# stage 1: install dependencies
FROM golang:1.18.0-alpine3.14 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# stage 2: build binary for production
FROM base as build

RUN go build -v -o main ./cmd/server

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

COPY --from=build /app/main .

ENTRYPOINT "./main"

# stage 3: run in watch mode
FROM base as dev

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT CompileDaemon --build="go build -v -o go-bin ./cmd/server/main.go" --command=./go-bin

