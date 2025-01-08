package main

import (
	"log"

	"github.com/turut4/social/internal/env"
)

func main() {
	config := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: config,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

}
