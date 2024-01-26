MIGRATE_PATH = db/migrations
POSTGRES_DB = db
POSTGRES_USER = postgres
POSTGRES_PASSWORD = Passw0rd
POSTGRES_HOST = 0.0.0.0
POSTGRES_PORT = 5432
POSTGRES_HOST_DOCKER = db

migrate-create:
	@echo "migrate create"
	@docker compose run --rm web sh -c "migrate -path /app/db/migration create -ext sql -dir db/migrations -seq $(NAME)"

migrate-up:
	@echo "migrate up"
	@docker compose run --rm web sh -c "migrate -path /app/db/migration \
	  -source file://${MIGRATE_PATH} -database  \
		'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST_DOCKER}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable' -verbose up $(VERSION)"

migrate-down:
	@echo "migrate down"
	@docker compose run --rm web sh -c "migrate -path /app/db/migration  \
	-source file://${MIGRATE_PATH}  \
	-database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST_DOCKER}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable' down $(VERSION)"

generate:
	@echo "sqlc"
	@docker compose run --rm web sh -c "CGO_ENABLED=1 sqlc generate"

# -ext	マイグレーションファイルの拡張子
# -dir	マイグレーションファイルを作成する場所
# -seq	マイグレーションファイルの名前