postgres:
	docker run --name go-bank --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it go-bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it go-bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

tidy:
	go mod tidy

vendor:
	go mod vendor

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/samirprakash/go-bank/db/sqlc Store

local:
	docker build -t gobank:latest .
	docker run --name gobank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgres://root:secret@go-bank:5432/simple_bank?sslmode=disable" gobank:latest

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock local