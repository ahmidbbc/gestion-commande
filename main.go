// Package main is the bootstrap: it wires dependencies and starts the app.
// It holds no business logic — that lives in internal/service (usecases).
package main

import (
	"log"

	"gestion-commande/internal/api"
	"gestion-commande/internal/service"
	"gestion-commande/internal/store"
)

func main() {
	repo := store.New()
	svc := service.New(repo)
	h := api.NewHandler(svc)
	log.Fatal(h.ListenAndServe(":8080"))
}
