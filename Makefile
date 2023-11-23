migratefilesup:
	migrate create -ext sql -dir db/migration -seq init_schema
	
# db/migration dir is in the project dir
# Copy the sql codes of "Inventory Management.sql" in db_info dir to -
# the migration file with "up" suffix in db/migration dir after you run  "$ make migratefilesup"
# add these 3 lines to a file with "down" suffix:

# DROP TABLE IF EXISTS categories;
# DROP TABLE IF EXISTS units;
# DROP TABLE IF EXISTS goods;

postgres:
	docker run --name postgresInv -p 5432:5432 -e POSTGRES_USER=mosleh -e POSTGRES_PASSWORD=1234 -d postgres:latest

postgresstop:
	docker stop postgresInv

postgresdown:
	docker rm postgresInv

createdb:
	docker exec -it postgresInv createdb --username=mosleh --owner mosleh inventory_management

dropdb:
	docker exec -it postgresInv dropdb --username=mosleh inventory_management

execdb: # access to database psql command line
	docker exec -it postgresInv psql -U mosleh -n inventory_management
	
migrateup: 
	migrate -path db/migration -database "postgres://mosleh:1234@localhost:5432/inventory_management?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://mosleh:1234@localhost:5432/inventory_management?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: migratefilesup postgres postgresstop postgresdown createdb dropdb execdb migrateup migratedown sqlc test server