FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

USER app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/internal/infra/postgres/migrations /app/internal/infra/postgres/migrations

EXPOSE 4000

CMD ["/app/server"]
