package main

import (
	"context"
	"log"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
	"github.com/DevSchmied/subscription-aggregation-service/internal/domain"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/handlers"
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/router"
	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"
	"github.com/google/uuid"
)

func main() {
	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Smoke test started")

	ctx := context.Background()

	// Initialize PostgreSQL connection pool
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Create subscription repository
	repo := postgres.NewSubscriptionRepo(pool)

	// Prepare test subscription data
	sub := domain.Subscription{
		ID:          uuid.New(),
		ServiceName: "Spotify Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
	}

	// Test create operation
	created, err := repo.Create(ctx, sub)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created: %+v\n", created)

	// Test get by ID
	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got: %+v\n", got)

	subH := handlers.NewSubscriptionsHandler(repo, 3*time.Second)

	rtr := router.NewRouter(router.Dependencies{
		Subscriptions: subH,
	})

	addr := "localhost"
	rtr.Run(addr + ":" + cfg.AppPort)

}
