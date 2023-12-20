TRANSACTIONS_DB_URL=postgresql://postgres:password@localhost:5432/transactions_db?sslmode=disable
ANALYTICS_DB_URL=postgresql://postgres:password@localhost:5432/analytics_db?sslmode=disable

prepare: postgres createdb migrateup seeddb

postgres:
	@echo "Creating a posrgres docker container"
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	@echo "Creating the database"
	docker exec -it postgres createdb --username=postgres --owner=postgres transactions_db
	docker exec -it postgres createdb --username=postgres --owner=postgres analytics_db


migrateup:
	migrate -path transactions/db/migration -database "$(TRANSACTIONS_DB_URL)" -verbose up
	migrate -path analytics/db/migration -database "$(ANALYTICS_DB_URL)" -verbose up

seeddb:
	@echo "Seeding the database"
	go run ./transactions/db/seeder

dropdb:
	@echo "Droping the database"
	docker exec -it postgres dropdb --username=postgres ecommerce_db
	docker exec -it postgres dropdb --username=postgres analytics_db

mockdb:
	mockgen -package mockdb -destination ./transactions/db/mock/store.go  ecommerce/transactions/db/sqlc Querier

mocktemp:
	mockgen -package mocktemporal -destination ./transactions/workflow/mock/temporal.go  go.temporal.io/sdk/client Client