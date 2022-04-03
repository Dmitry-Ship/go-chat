# go-chat [![Front](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml) [![Back](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml)

real time chat app

[https://coverfield.app](https://coverfield.app)

## ⚡️ Quick Start

1. Install and boot up Docker
2. Create `.env` file at root directory taking `.env.example` as a base.
3. Run the following command:

```
docker-compose up --build
```

4. open browser and go to http://localhost:8080

## ☁️ Hosting platforms

- backend: https://console.cloud.google.com/run/detail/us-central1/go-chat/metrics?project=go-playground-311723
- frontend: https://vercel.com/dmitry-ship/go-chat

## ⚙️ Architecture overview

```mermaid
graph LR
    A(NGINX) --> B(NextJS)
    A --> C(Golang)
    A --> D(Golang)
    C --> E(Postgres)
    D --> E
    C --> F(Redis PubSub)
    D --> F
```
