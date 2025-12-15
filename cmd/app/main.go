package main

import (
	"fmt"
	"log"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
)

func main() {
	// Load application config
	cfg, err := config.Load()
	_ = cfg
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello subscription-aggregation-service")
}
