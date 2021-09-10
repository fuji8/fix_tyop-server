FROM golang:1.16.7-alpine as server-build

WORKDIR /github.com/fuji8/fix_tyop-server

COPY go.mod go.sum ./
ENV GO111MODULE=on
RUN go mod download
COPY ./ ./

RUN go build -o main

FROM alpine:3.13.5

WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1

RUN apk --update add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \ 
    && apk add --update ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

COPY --from=server-build /github.com/fuji8/fix_tyop-server/main ./

ENTRYPOINT ./main
