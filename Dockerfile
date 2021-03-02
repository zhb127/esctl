################################
# builder
################################
FROM golang:1.15-alpine AS builder

RUN set -ex \
  && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk update \
  && apk add --no-cache alpine-sdk git \
  && rm -rf /var/cache/apk/*

RUN mkdir -p /build
WORKDIR /build

COPY . .
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE="auto"

RUN go mod download

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -o /build/esctl ./main.go

################################
# runtime
################################
FROM alpine:3.11

RUN  set -ex \
  && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk update \
  && apk add --no-cache ca-certificates tzdata bash \
  && rm -rf /var/cache/apk/*

ENV TZ=Asia/Shanghai

COPY --from=builder /build/* /usr/local/bin/
