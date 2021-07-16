FROM golang:1.13 as builder
ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV GOPROXY=https://goproxy.cn

WORKDIR /go/src/goodrain.com/operatelog
COPY . .

ARG GOARCH
ARG LDFLAGS
RUN go build -ldflags "$LDFLAGS" -o /log $PWD/cmd/log

FROM alpine:3.11.2
RUN apk add --update tzdata \
    && mkdir /app \
    && apk add --update apache2-utils \
    && rm -rf /var/cache/apk/*
ENV TZ=Asia/Shanghai
COPY --from=builder log /app

ENTRYPOINT ["/app/log"]