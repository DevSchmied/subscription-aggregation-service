package handlers

import (
	"context"
	"errors"
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

// parseSubscriptionRequest validates input and builds domain entity.
func parseSubscriptionRequest(
	req SubscriptionRequest,
	id uuid.UUID,
) (domain.Subscription, error) {

	serviceName := strings.TrimSpace(req.ServiceName)
	if serviceName == "" {
		return domain.Subscription{}, errors.New("service_name is required")
	}

	if req.Price < 0 {
		return domain.Subscription{}, errors.New("price must be >= 0")
	}

	userID, err := uuid.Parse(strings.TrimSpace(req.UserID))
	if err != nil {
		return domain.Subscription{}, errors.New("invalid user_id")
	}

	start, err := utils.ParseMonthYear(strings.TrimSpace(req.StartDate))
	if err != nil {
		return domain.Subscription{}, err
	}

	// Parse optional end date
	var endPtr *time.Time
	if end := strings.TrimSpace(req.EndDate); end != "" {
		parsedEnd, err := utils.ParseMonthYear(end)
		if err != nil {
			return domain.Subscription{}, err
		}
		if parsedEnd.Before(start) {
			return domain.Subscription{}, errors.New("end_date before start_date")
		}
		endPtr = &parsedEnd
	}

	// Build domain model
	return domain.Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     endPtr,
	}, nil
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

	// Parse request
	sub, err := parseSubscriptionRequest(req, uuid.New())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

// Update handles subscription update request.
func (h *SubscriptionsHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// Parse request
	sub, err := parseSubscriptionRequest(req, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.dbTimeout)
	defer cancel()

	out, err := h.repo.Update(ctx, sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, toResponse(out))
}

// Delete removes subscription by ID.
func (h *SubscriptionsHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.dbTimeout)
	defer cancel()

	if err := h.repo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
