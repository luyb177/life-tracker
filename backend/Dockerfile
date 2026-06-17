FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 主服务
RUN go build -ldflags="-s -w" -o /app/life_tracker life_tracker.go

# 定时任务
RUN go build -ldflags="-s -w" -o /app/cron ./cmd/cron/


FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/life_tracker /app/cron /usr/local/bin/

EXPOSE 8888

CMD ["life_tracker", "-f", "etc/life_tracker.yaml"]
