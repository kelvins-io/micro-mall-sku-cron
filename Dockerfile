FROM golang:1.13.15-buster as gobuild
WORKDIR $GOPATH/src/gitee.com/cristiane/micro-mall-sku-cron
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN bash ./build.sh
## FROM alpine:latest as gorun
FROM ubuntu:latest as gorun
WORKDIR /www/
COPY --from=gobuild /go/src/gitee.com/cristiane/micro-mall-sku-cron/micro-mall-sku-cron .
COPY --from=gobuild /go/src/gitee.com/cristiane/micro-mall-sku-cron/etc ./etc
CMD ["./micro-mall-sku-cron"]
