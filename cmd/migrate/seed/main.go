package main

import (
	"log"

	"github.com/turut4/social/internal/db"
	"github.com/turut4/social/internal/env"
	"github.com/turut4/social/internal/store"
)

func main() {
	addr := env.GetString("DB _ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
