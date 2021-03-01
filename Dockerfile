FROM golang:1.16-alpine AS builder

ARG bin
ARG dir

WORKDIR $dir

COPY . $dir

RUN apk update && apk add --update alpine-sdk
RUN make NAME=${bin}
