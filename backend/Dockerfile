# stage 1: install dependencies
FROM golang:1.18.3-alpine3.16 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . .

# stage 2: build binary for production
FROM base as build

RUN go build -v -o main ./cmd/server

# stage 3: run binary
FROM alpine:3.16 as prod

RUN apk --no-cache add ca-certificates

COPY --from=build /app/main .

ENTRYPOINT "./main"

# stage 4: run in watch mode
FROM base as dev

ENTRYPOINT CompileDaemon --build="go build -v -o go-bin ./cmd/server/main.go" --command=./go-bin
