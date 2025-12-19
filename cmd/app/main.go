package main

import (
	"context"
	"log"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/handlers"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/router"
	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"
)

func main() {
	// Load application config
	cfg, err := config.Load()
	_ = cfg
	if err != nil {
		log.Fatal(err)
	}

	log.Println("App started")

	ctx := context.Background()

	// Initialize PostgreSQL connection pool
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Create subscription repository
	repo := postgres.NewSubscriptionRepo(pool)

	subH := handlers.NewSubscriptionsHandler(repo, 3*time.Second)

	rtr := router.NewRouter(router.Dependencies{
		Subscriptions: subH,
	})
	addr := "localhost"
	rtr.Run(addr + ":" + cfg.AppPort)
}
