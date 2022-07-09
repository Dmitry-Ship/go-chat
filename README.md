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

## 💿 Database Design

https://drawsql.app/none-794/diagrams/go-chat/embed

## 🌄 Screenshots

<img width="260" alt="Screenshot 2022-06-04 at 21 38 35" src="https://user-images.githubusercontent.com/21371972/172021307-20440dbb-215f-4339-8c70-cfcb2fe8bb4a.png"> <img width="260" alt="Screenshot 2022-06-04 at 21 37 05" src="https://user-images.githubusercontent.com/21371972/172021316-fbfc2534-7934-460d-9bce-48cb4174c25f.png"> <img width="260" alt="Screenshot 2022-06-04 at 21 40 54" src="https://user-images.githubusercontent.com/21371972/172021326-435029f9-09ea-476e-bdbb-25838a2b697f.png">

<img width="260" alt="Screenshot 2022-06-04 at 21 36 46" src="https://user-images.githubusercontent.com/21371972/172021335-ec9efe67-de77-4996-bfc2-d42652f0383e.png"> <img width="260" alt="Screenshot 2022-06-04 at 21 36 12" src="https://user-images.githubusercontent.com/21371972/172021338-c54633c3-b49d-4163-8110-6db62c16281c.png"> <img width="260" alt="Screenshot 2022-06-04 at 21 39 05" src="https://user-images.githubusercontent.com/21371972/172021343-476f2dde-2461-4488-83c8-760baf393968.png">
