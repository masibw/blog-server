FROM golang:1.15.6 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/github.com/masibw/blog-server

COPY go.mod  ./
RUN go mod download

COPY . .
RUN go build .

FROM alpine:3.13.0

ENV DOCKERIZE_VERSION v0.6.1

RUN apk add --no-cache bash ca-certificates curl openssl

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=build /go/src/github.com/masibw/blog-server/blog-server /

RUN chmod a+x /blog-server