package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/domain"
	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"
	"github.com/DevSchmied/subscription-aggregation-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionsHandler handles subscription HTTP requests.
type SubscriptionsHandler struct {
	repo      *postgres.SubscriptionRepo
	dbTimeout time.Duration
}

// NewSubscriptionsHandler creates a new subscriptions handler.
func NewSubscriptionsHandler(
	repo *postgres.SubscriptionRepo,
	timeout time.Duration,
) *SubscriptionsHandler {
	return &SubscriptionsHandler{
		repo:      repo,
		dbTimeout: timeout,
	}
}

// Create handles subscription creation request.
func (h *SubscriptionsHandler) Create(c *gin.Context) {
	body := make(map[string]string)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	// Normalize user input
	serviceName := strings.TrimSpace(body["service_name"])
	priceStr := strings.TrimSpace(body["price"])
	userIDStr := strings.TrimSpace(body["user_id"])
	startDateStr := strings.TrimSpace(body["start_date"])
	endDateStr := strings.TrimSpace(body["end_date"])

	// Validate required fields
	if serviceName == "" || priceStr == "" || userIDStr == "" || startDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required fields",
		})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid price",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	start, err := utils.ParseMonthYear(startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Parse optional end date
	var endPtr *time.Time
	if endDateStr != "" {
		end, err := utils.ParseMonthYear(endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if end.Before(start) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "end_date before start_date",
			})
			return
		}
		endPtr = &end
	}

	// Build domain model
	sub := domain.Subscription{
		ID:          uuid.New(),
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     endPtr,
	}

	// Apply database timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.dbTimeout)
	defer cancel()

	out, err := h.repo.Create(ctx, sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}

	// Format optional end date
	endDate := func() string {
		if out.EndDate != nil {
			return utils.FormatMonthYear(*out.EndDate)
		}
		return ""
	}()

	// Return created subscription
	c.JSON(http.StatusCreated, gin.H{
		"id":           out.ID.String(),
		"service_name": out.ServiceName,
		"price":        out.Price,
		"user_id":      out.UserID.String(),
		"start_date":   utils.FormatMonthYear(out.StartDate),
		"end_date":     endDate,
		"created_at":   out.CreatedAt.Format(time.RFC3339),
		"updated_at":   out.UpdatedAt.Format(time.RFC3339),
	})
}
