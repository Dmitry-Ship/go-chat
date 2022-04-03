FROM node:16-alpine3.14 AS base

WORKDIR /client

COPY package*.json ./

RUN npm install

COPY . .

FROM base as build

RUN npm run build

FROM node:16-alpine3.14 AS prod

WORKDIR /client

ENV NODE_ENV production

COPY --from=build /app/next.config.js ./
COPY --from=build /client/public ./public
COPY --from=build /client/package.json ./package.json
COPY --from=build --chown=nextjs:nodejs /client/.next/standalone ./
COPY --from=build --chown=nextjs:nodejs /client/.next/static ./.next/static

EXPOSE 3000

CMD ["node", "server.js"]

FROM base as dev

EXPOSE 3000

CMD [ "npm", "run" , "dev"]

