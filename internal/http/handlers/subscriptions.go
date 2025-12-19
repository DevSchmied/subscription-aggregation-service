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

// toResponse maps domain subscription to API response.
func toResponse(s domain.Subscription) SubscriptionResponse {
	end := ""
	// Format optional end date
	if s.EndDate != nil {
		end = utils.FormatMonthYear(*s.EndDate)
	}

	// Build API response
	return SubscriptionResponse{
		ID:          s.ID.String(),
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   utils.FormatMonthYear(s.StartDate),
		EndDate:     end,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339), // ISO timestamp format
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}

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

	resp := toResponse(out)

	// Return created subscription
	c.JSON(http.StatusCreated, resp)
}

// Get returns subscription by ID.
func (h *SubscriptionsHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.dbTimeout) // limit db time
	defer cancel()

	s, err := h.repo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, toResponse(s))
}
