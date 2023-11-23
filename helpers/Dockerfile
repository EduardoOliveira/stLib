FROM golang:1.19 AS go

RUN mkdir app
WORKDIR /app
COPY . .

RUN go build main.go
RUN ls -la helpers
RUN pwd

FROM debian:latest
RUN mkdir /config

WORKDIR /app

COPY --from=go /app/main ./
COPY --from=go /app/helpers/config.toml /config

RUN mkdir /library
RUN mkdir /app/cache
RUN ls -la /config
RUN pwd
RUN chmod +x main
RUN ls -la /app/main

VOLUME [ "/config", "/library" ]
ENTRYPOINT ["/app/main"]