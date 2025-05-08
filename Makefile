migration_path := internal/storage/postgres/migrations
migrateinit:
	migrate create -ext sql -dir $(migration_path) -seq init_schema
migrateup:
	migrate -path $(migration_path) -database "postgresql://root:secret@localhost:5433/postgres?sslmode=disable" -verbose up
migratedown:
	migrate -path $(migration_path) -database "postgresql://root:secret@localhost:5433/postgres?sslmode=disable" -verbose down
connectdb:
	sudo docker exec -it pg_warehouse psql -U root postgres
createcontainer:
	sudo docker run --name pg_warehouse -p 5433:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d f2258b53bc9c
sqlc:
	sqlc generate
test:
	go test ./internal/rest

.PHONY: connectdb createcontainer migrateup migratedown migrateinit sqlc test