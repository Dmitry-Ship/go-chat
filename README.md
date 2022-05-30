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
    A(User) --> B(Nginx)
    B --> C(NextJS)
    B --> D(Golang)
    B --> E(Golang)
    E --> F(Postgres)
    E --> F
    D --> G(Redis PubSub)
    E --> G
```

## 💿 Database Design

https://drawsql.app/none-794/diagrams/go-chat/embed

## 🌄 Screenshots

<img width="260" alt="Screenshot 2022-05-29 at 16 36 42" src="https://user-images.githubusercontent.com/21371972/170965594-3d9db99b-a3b6-4ff7-9d19-73cd8029b4ad.png"> <img width="260" alt="Screenshot 2022-05-29 at 16 29 42" src="https://user-images.githubusercontent.com/21371972/170965835-716d7f8a-30de-40a2-b5ff-0b97fd4b8007.png"> <img width="260" alt="Screenshot 2022-05-30 at 12 45 19" src="https://user-images.githubusercontent.com/21371972/170965955-f586fcb9-0efb-46a7-9c58-9e4d6501f317.png">


<img width="260" alt="Screenshot 2022-05-29 at 16 36 55" src="https://user-images.githubusercontent.com/21371972/170966050-34ad04bf-d115-4505-8b84-c5c3a0255a26.png"> <img width="260" alt="Screenshot 2022-05-29 at 16 28 08" src="https://user-images.githubusercontent.com/21371972/170966086-4f4af6d5-9892-4453-9f2b-46e38ba50e0d.png"> <img width="260" alt="Screenshot 2022-05-29 at 16 27 43" src="https://user-images.githubusercontent.com/21371972/170966128-50bea454-000c-43bb-a50d-14b9ae926e05.png">

 
