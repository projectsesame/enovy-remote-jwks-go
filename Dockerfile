FROM --platform=$BUILDPLATFORM golang:1.19 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum
COPY jwks.json jwks.json

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main main.go

FROM docker.m.daocloud.io/alpine:3.16.2

COPY --from=builder /app/main /bin/
COPY --from=builder /app/jwks.json /jwks.json


EXPOSE 8080
