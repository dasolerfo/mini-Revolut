DB_URL=postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable

postgres: 
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Songoku13 -d postgres:12-alpine
startDB:
	docker start postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateupremote:
	migrate -path db/migration -database "postgresql://root:5KmWKrX7oKF7rEUckJwK@bank-system.cfg8weceu069.eu-west-3.rds.amazonaws.com:5432/bank_system" -verbose up

migratedownremote:
	migrate -path db/migration -database "postgresql://root:5KmWKrX7oKF7rEUckJwK@bank-system.cfg8weceu069.eu-west-3.rds.amazonaws.com:5432/bank_system" -verbose down


migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

upgradesqlc:
	brew upgrade sqlc
db_docs:
	 dbdocs build docs/db.dbml
db_schema:
	 dbml2sql --postgres -o docs/schema.sql docs/db.dbml	 

test: 
	go test -v -cover ./...

mockgen:
	mockgen -package mockdb -destination db/mock/store.go ./db/model Store

serve: 
	go run main.go

proto: 
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
	
evans:
	evans --host localhost --port 8081 -r repl
	
.PHONY: createdb startDB dropdb postgres migrateup migratedown sqlc test serve upgradesqlc migrateupremote migratedownremote migrateup1 migratedown1 mockgen db_docs db_schemas proto


