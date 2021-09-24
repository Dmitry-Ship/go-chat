# stage 1: install dependencies
FROM golang:1.16.7-alpine AS base

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

# stage 4: install dev dependencies 
FROM base as dev

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -v -o go-bin ./cmd/server/main.go" --command=./go-bin

