### ビルド用ステージ ###
FROM golang:1.21.4-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
# CGoを無効化してビルド
# RUN GOOS=linux GOARCH=amd64 go build -o main .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

### 開発用ステージ ###
FROM golang:1.21.4-alpine as development

# アップデートとvim,gitのインストール
RUN apk update && apk add vim git

# gccとpostgresqlのインストール。sqlcコマンドで必要
RUN apk add --no-cache gcc musl-dev postgresql-dev

# sqlcのインストール
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# migrateのインストール
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# /appに移動
WORKDIR /app

# builderからmainをコピー
COPY --from=builder /app/main .

# airのインストール
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]

### 本番用ステージ###
FROM alpine:latest as production
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
