# go
go mod init twitter_golang_backend
go mod tidy
go get -u github.com/gin-gonic/gin
go get -u golang.org/x/crypto/bcrypt
go get github.com/lib/pq
go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

# docker compose
## docker composeの操作方法
docker compose build
docker compose up -d

## prune で不要なコンテナ、ネットワーク、ボリューム、イメージを削除
docker system prune
docker system prune -a

## コンテナ、ネットワーク、ボリューム、イメージを削除
docker compose down
docker compose down -v
docker compose down --rmi all --volumes --remove-orphans

# Makefile
## マイグレーションファイルの作成
make migrate-create NAME={ TABLE_NAME }
make migrate-create NAME=create_users2
## マイグレーションの実行(VERSIONを指定するとそのバージョンまで実行する)
make migrate-up
make migrate-up VERSION="1"
## マイグレーションのロールバック(VERSIONを指定するとそのバージョンまでロールバックする)
make migrate-down
make migrate-down VERSION=1
## sqlcの実行
make generate

# migrate
## マイグレーションを実行
docker compose run --rm web migrate -path /app/db/migration create -ext sql -dir db/migrations -seq users
docker compose run --rm web migrate -path /app/db/migration -database "postgres://postgres:Passw0rd@db:5432/db?sslmode=disable" up
docker compose run --rm web migrate -path /app/db/migration -database "postgres://postgres:Passw0rd@db:5432/db?sslmode=disable" down

## データベースへの接続
docker compose exec db psql -U postgres -d db
\dt
\d users

## マイグレーション生成
### sqlcを実行することで、migrationファイルからGoのコードを生成する
docker compose run --rm -e CGO_ENABLED=1 web sqlc generate
docker-compose exec web ls /app/db/migration

# Redis
## go get
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/sessions/redis
go get github.com/google/uuid

## redis-cli
docker-compose exec redis redis-cli
keys *
get "key"
del "key"
