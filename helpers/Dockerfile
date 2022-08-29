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

FROM debian:latest
RUN mkdir -p app/frontend/dist
WORKDIR /app

COPY --from=go /go/app/main ./
COPY --from=go /go/app/helpers/config.toml ./
COPY --from=n /app/frontend/dist ./frontend/
RUN mkdir /library
RUN ls -la
RUN pwd
RUN chmod +x main
RUN ls -la /app/main
ENTRYPOINT ["/app/main"]
