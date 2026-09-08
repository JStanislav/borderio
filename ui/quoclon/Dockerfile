# syntax=docker/dockerfile:1
# Build the app
FROM node:24-alpine3.23 AS builder

WORKDIR /app

ARG VITE_SERVER_URL
ENV VITE_SERVER_URL=${VITE_SERVER_URL}

ARG APP_PORT
ENV VITE_APP_PORT=${APP_PORT}

COPY package*.json ./

RUN npm i

COPY --exclude=node_modules . .

RUN npm run build


# Runner 
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

ENV NODE_ENV=production

CMD ["nginx", "-g", "daemon off;"]