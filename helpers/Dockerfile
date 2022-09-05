FROM node:18-alpine3.15 AS n

RUN mkdir app
WORKDIR /app
COPY . .
WORKDIR /app/frontend
RUN pwd
RUN ls -la
RUN npm install
RUN npm run build

FROM golang:1.19 AS go

RUN mkdir app
WORKDIR /app
COPY . .

RUN go build main.go
RUN ls -la helpers
RUN pwd

FROM debian:latest
RUN mkdir -p app/frontend/dist
RUN mkdir /config

WORKDIR /app

COPY --from=go /app/main ./
COPY --from=go /app/helpers/config.toml /config
COPY --from=n /app/frontend/dist ./frontend/dist
RUN mkdir /library
RUN mkdir /app/cache
RUN ls -la /config
RUN pwd
RUN chmod +x main
RUN ls -la /app/main

ENV STLIB_CONFIG=/config/config.toml

VOLUME [ "/config", "/library" ]
ENTRYPOINT ["/app/main"]