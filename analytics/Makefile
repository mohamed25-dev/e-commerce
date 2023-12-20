# DB_URL=postgresql://postgres:password@localhost:5432/analytics_db?sslmode=disable
DB_URL=cockroachdb://mohamed:@localhost:9999/analytics_db?sslmode=disable

createdb:
	@echo "Creating the database"
	docker exec -it postgres createdb --username=postgres --owner=postgres analytics_db

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

seeddb:
	@echo "Seeding the database"
	go run ./transactions/db/seeder

dropdb:
	@echo "Droping the database"
	docker exec -it postgres dropdb --username=postgres analytics_db


sqlc:
	@echo "Generating sqlc code"
	sqlc generate

protoc:
	@echo "Generating proto code"
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./proto/analytics.proto

serve:
	@echo "Running the server"
	go run .