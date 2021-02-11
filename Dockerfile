FROM golang:alpine3.13

WORKDIR /go

COPY ./src ./jato/src

RUN set -ex \
    && apk update \
    && apk add --no-cache \
    git \
    && go get golang.org/x/crypto/ssh