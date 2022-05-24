# go-chat [![Front](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/front.yml) [![Back](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml/badge.svg)](https://github.com/Dmitry-Ship/go-chat/actions/workflows/back.yml)

Real time multi room chat app, built solely for educational purposes.

## ⚡️ Quick Start

1. Install and boot up Docker
2. Create `.env` file at the root directory taking `.env.example` as a base.
3. Run the following command:

```
docker-compose up --build
```

4. Go to http://localhost:8080

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

## 💿 Database Design

https://drawsql.app/none-794/diagrams/go-chat/embed

## 🌄 Screenshots
<img width="260" alt="Screenshot 2022-05-24 at 19 23 34" src="https://user-images.githubusercontent.com/21371972/170085483-a3fce839-6e22-422e-8d38-90bef88ca716.png"> <img width="260" alt="Screenshot 2022-05-24 at 19 23 43" src="https://user-images.githubusercontent.com/21371972/170085493-016ce553-ff56-4f22-950b-72347227e36c.png"> <img width="260" alt="Screenshot 2022-05-24 at 19 23 57" src="https://user-images.githubusercontent.com/21371972/170085504-0b61723c-50ff-4839-bf16-d23116480178.png">
