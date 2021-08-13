FROM node:14.17.0-alpine AS base

WORKDIR /client

COPY package*.json ./

RUN npm install

COPY . .

FROM base as build

RUN go build -v -o main ./cmd/server

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

COPY --from=build /client/build .

FROM base as dev

EXPOSE 3000

CMD [ "npm", "run" , "dev"]

