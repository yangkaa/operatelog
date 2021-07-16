FROM golang:1.13 as builder
ADD . /go/src/github.com/operatelog
WORKDIR /go/src/github.com/operatelog
ENV GOPROXY=https://goproxy.cn

RUN go build -ldflags "-w -s" -o _output/linux-amd64/operatelog ./cmd/operatelog

FROM node:14 as uibuilder
ADD ui /ui
WORKDIR /ui
RUN npm config set registry https://registry.npm.taobao.org && npm install && npm run build

FROM goodrainapps/alpine:3.4
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /go/src/github.com/operatelog/_output/linux-amd64/operatelog /app/operatelog
COPY --from=uibuilder /ui/dist /app/ui/dist
COPY --from=uibuilder /ui/web.conf /etc/nginx/conf.d/web.conf
WORKDIR /app
ENV PORT=8080
EXPOSE 8080
CMD ["./operatelog"]
