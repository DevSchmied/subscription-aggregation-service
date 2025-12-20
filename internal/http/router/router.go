package router

import (
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/handlers"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Dependencies groups all HTTP handlers required by the router.
type Dependencies struct {
	Subscriptions *handlers.SubscriptionsHandler
	Aggregation   *handlers.AggregationHandler
}

// NewRouter configures and returns a Gin HTTP router.
func NewRouter(d Dependencies) *gin.Engine {
	// Create Gin engine without default middleware
	r := gin.New()

	// Register logging and panic recovery middleware
	r.Use(gin.Logger(), gin.Recovery())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		// CRUDL operations for subscriptions
		api.POST("/subscriptions", d.Subscriptions.Create)
		api.GET("/subscriptions/:id", d.Subscriptions.Get)
		api.PUT("/subscriptions/:id", d.Subscriptions.Update)
		api.DELETE("/subscriptions/:id", d.Subscriptions.Delete)
		api.GET("/subscriptions", d.Subscriptions.List)

		// Aggregation endpoint: calculate total subscription cost for a period
		api.GET("/subscriptions/total", d.Aggregation.Total)
	}

	return r
}
