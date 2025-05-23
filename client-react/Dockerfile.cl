FROM node:24.0-bullseye AS build

WORKDIR /app

COPY package*.json ./

RUN --mount=type=cache,target=/app/.npm \
    npm set cache /app/npm && \
    npm install

COPY . .

RUN npm run build 

FROM nginxinc/nginx-unprivileged:1.23-alpine-perl

COPY --link nginx.conf /etc/nginx/conf.d/default.conf

COPY --link --from=build /app/dist/ /usr/share/nginx/html

EXPOSE 8080