# go-chat [![Front](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml) [![Back](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml)

Real time multi room chat app, built solely for educational purposes.

## ⚡️ Quick Start

1. Install and boot up Docker
2. Run `cp .env.example .env` and tweak it to your needs
3. Run `docker-compose up --build`
4. Go to http://localhost:8080

## ⚙️ Architecture overview

```mermaid
graph LR
    A(User) --> B(Nginx)
    B --> C(NextJS)
    B --> D(Golang)
    B --> E(Golang)
    D --> F(Postgres)
    E --> F
    D --> G(Redis PubSub)
    E --> G
```
