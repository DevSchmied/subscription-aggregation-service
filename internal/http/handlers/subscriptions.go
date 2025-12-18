package handlers

import (
	"context"
	"net/http"
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

// SubscriptionRequest defines create payload.
type SubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // MM-YYYY
	EndDate     string `json:"end_date"`
}

// SubscriptionResponse defines API response.
type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Create handles subscription creation request.
func (h *SubscriptionsHandler) Create(c *gin.Context) {
	req := SubscriptionRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	serviceName := strings.TrimSpace(req.ServiceName)
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_name is required"})
		return
	}

	if req.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be >= 0"})
		return
	}

	userID, err := uuid.Parse(strings.TrimSpace(req.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	start, err := utils.ParseMonthYear(strings.TrimSpace(req.StartDate))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse optional end date
	var endPtr *time.Time
	if end := strings.TrimSpace(req.EndDate); end != "" {
		parsedEnd, err := utils.ParseMonthYear(end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if parsedEnd.Before(start) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date before start_date"})
			return
		}
		endPtr = &parsedEnd
	}

	// Build domain entity
	sub := domain.Subscription{
		ID:          uuid.New(),
		ServiceName: serviceName,
		Price:       req.Price,
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

	// Build API response
	resp := SubscriptionResponse{
		ID:          out.ID.String(),
		ServiceName: out.ServiceName,
		Price:       out.Price,
		UserID:      out.UserID.String(),
		StartDate:   utils.FormatMonthYear(out.StartDate),
		EndDate:     endDate,
		CreatedAt:   out.CreatedAt.Format(time.RFC3339), // ISO timestamp format
		UpdatedAt:   out.UpdatedAt.Format(time.RFC3339),
	}

	// Return created subscription
	c.JSON(http.StatusCreated, resp)
}
