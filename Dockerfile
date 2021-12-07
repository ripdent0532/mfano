FROM golang:alpine as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
RUN apk add gcc && apk add g++
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add tzdata ca-certificates && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata && rm -rf /var/cache/apk/*
ENV LANG="C.UTF-8" \
    HOME_DIR=/home \
    DST_DIR=/root/pp/static \
    DB_DIR=/root/pp/db
WORKDIR /root
COPY --from=builder /app/mfano .
COPY mfano.sqlite /root
COPY views views
ADD entrypoint.sh .
RUN chomd +x entrypoint.sh

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
ENTRYPOINT ["/root/entrypoint.sh"]
CMD ["/root/mfano"]