// @title Subscription Aggregation API
// @version 1.0
// @description REST service for aggregating user subscription costs
// @host localhost:8080
// @BasePath /api
package main

import (
	"context"
	"log"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/handlers"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/router"
	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"

	_ "github.com/DevSchmied/subscription-aggregation-service/docs"
)

func main() {
	// Load application configuration (env)
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

	// Initialize HTTP handlers with DB timeout
	subH := handlers.NewSubscriptionsHandler(repo, 3*time.Second)
	aggH := handlers.NewAggregationHandler(repo, 3*time.Second)

	// Build HTTP router and inject dependencies
	rtr := router.NewRouter(router.Dependencies{
		Subscriptions: subH,
		Aggregation:   aggH,
	})

	// Start HTTP server
	log.Println("HTTP server started on port", cfg.AppPort)
	if err := rtr.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
