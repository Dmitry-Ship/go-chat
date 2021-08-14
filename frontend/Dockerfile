FROM node:14.17.0-alpine AS base

WORKDIR /client

COPY package*.json ./

RUN npm install

COPY . .

FROM base as build

RUN npm run build

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

COPY --from=build /client/dist .

FROM base as dev

EXPOSE 3000

CMD [ "npm", "run" , "dev"]
