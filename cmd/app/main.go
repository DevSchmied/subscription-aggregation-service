package main

import (
	"log"
	"net/http"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
)

func main() {
	// Load application config
	cfg, err := config.Load()
	_ = cfg
	if err != nil {
		log.Fatal(err)
	}

	log.Println("App started")
	http.ListenAndServe(":8080", nil)
}
