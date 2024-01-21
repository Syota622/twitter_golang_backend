### ビルド用ステージ ###
FROM golang:1.21.4-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
# CGoを無効化してビルド
RUN GOOS=linux GOARCH=amd64 go build -o main .

### 開発用ステージ ###
FROM golang:1.21.4-alpine as development

# アップデートとgitのインストール
RUN apk update && apk add vim git

# sqlcのインストール
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.22.0

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
