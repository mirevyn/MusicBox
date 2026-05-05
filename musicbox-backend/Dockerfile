# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o musicbox-server main.go


# 运行阶段
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/musicbox-server .
COPY --from=builder /app/config ./config

RUN mkdir -p uploads/music uploads/covers uploads/lyrics uploads/avatars

EXPOSE 8000

CMD ["./musicbox-server"]
