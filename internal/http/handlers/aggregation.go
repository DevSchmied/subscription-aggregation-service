package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"
	"github.com/DevSchmied/subscription-aggregation-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AggregationHandler handles aggregation HTTP endpoints.
type AggregationHandler struct {
	repo      *postgres.SubscriptionRepo
	dbTimeout time.Duration
}

// NewAggregationHandler creates aggregation handler.
func NewAggregationHandler(repo *postgres.SubscriptionRepo, dbTimeout time.Duration) *AggregationHandler {
	return &AggregationHandler{
		repo:      repo,
		dbTimeout: dbTimeout,
	}
}

// Total calculates subscription cost for a given period.
// The sum includes only months when subscriptions were active.
func (h *AggregationHandler) Total(c *gin.Context) {
	startStr := strings.TrimSpace(c.Query("start_date"))
	endStr := strings.TrimSpace(c.Query("end_date"))

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date required",
		})
		return
	}

	periodStart, err := utils.ParseMonthYear(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date",
		})
		return
	}

	periodEnd, err := utils.ParseMonthYear(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date",
		})
		return
	}

	// Optional user filter
	var userID *uuid.UUID
	if v := strings.TrimSpace(c.Query("user_id")); v != "" {
		parsedID, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user_id",
			})
			return
		}
		userID = &parsedID
	}

	// Optional service filter
	var serviceName *string
	if v := strings.TrimSpace(c.Query("service_name")); v != "" {
		serviceName = &v
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.dbTimeout)
	defer cancel()

	// Fetch overlapping subscriptions
	items, err := h.repo.ListOverlapping(
		ctx,
		userID,
		serviceName,
		periodStart,
		periodEnd,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}

	total := 0

	for _, s := range items {

		// Determine actual start date of subscription
		// within the requested period
		activeStart := periodStart
		if s.StartDate.After(periodStart) {
			activeStart = s.StartDate
		}

		// Determine actual end date of subscription
		// within the requested period
		activeEnd := periodEnd
		if s.EndDate != nil && s.EndDate.Before(periodEnd) {
			activeEnd = *s.EndDate
		}

		// Skip subscriptions not active
		// during the requested period
		if activeEnd.Before(activeStart) {
			continue
		}

		// Calculate number of active months (inclusive)
		months := (activeEnd.Year()-activeStart.Year())*12 +
			int(activeEnd.Month()-activeStart.Month()) + 1

		// Add subscription cost for active months
		total += months * s.Price
	}

	// Return aggregation result
	c.JSON(http.StatusOK, gin.H{
		"total":        total,
		"period_start": startStr,
		"period_end":   endStr,
	})
}
