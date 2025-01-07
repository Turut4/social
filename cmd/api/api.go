package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr        string
	rateLimiter string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /health", app.healthCheckHandler)
	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("server have started at %s", app.config.addr)
	return srv.ListenAndServe()
}
