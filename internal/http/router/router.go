package router

import (
	"github.com/DevSchmied/subscription-aggregation-service/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Subscriptions *handlers.SubscriptionsHandler
}

func NewRouter(d Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/subscriptions", d.Subscriptions.Create)

	return r
}
