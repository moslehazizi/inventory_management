package main

import (
	"database/sql"
	"inventory_management/api"
	db "inventory_management/db/sqlc"
	_ "github.com/lib/pq"

	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://mosleh:1234@localhost:5432/inventory_management?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Connot connect to the database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("connot start server:", err)
	}

}
