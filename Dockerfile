FROM alpine:latest

MAINTAINER ServerScan

RUN mkdir -p /app

WORKDIR /app

ADD ./ServerScan /app

ENTRYPOINT ["./ServerScan"]
