package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/domain"
	"github.com/DevSchmied/subscription-aggregation-service/internal/storage/postgres"
	"github.com/google/uuid"
)

// ==============================================================
// ==============================================================
// toResponse
// ==============================================================
// ==============================================================
func TestToResponse_WithoutEndDate(t *testing.T) {
	// Arrange
	startDate := time.Date(2025, 7, 15, 10, 0, 0, 0, time.UTC)
	createdAt := time.Date(2025, 7, 1, 9, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2025, 7, 2, 11, 45, 0, 0, time.UTC)

	sub := domain.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   startDate,
		EndDate:     nil,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Act
	resp := toResponse(sub)

	// Assert

	// StartDate -> MM-YYYY
	if resp.StartDate != "07-2025" {
		t.Errorf("expected StartDate 07-2025, got %s", resp.StartDate)
	}

	// EndDate -> empty string
	if resp.EndDate != "" {
		t.Errorf("expected empty EndDate, got %s", resp.EndDate)
	}

	// CreatedAt -> RFC3339
	expectedCreated := createdAt.Format(time.RFC3339)
	if resp.CreatedAt != expectedCreated {
		t.Errorf("expected CreatedAt %s, got %s", expectedCreated, resp.CreatedAt)
	}

	// UpdatedAt -> RFC3339
	expectedUpdated := updatedAt.Format(time.RFC3339)
	if resp.UpdatedAt != expectedUpdated {
		t.Errorf("expected UpdatedAt %s, got %s", expectedUpdated, resp.UpdatedAt)
	}
}

func TestToResponse_WithEndDate(t *testing.T) {
	// Arrange
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 20, 0, 0, 0, 0, time.UTC)

	sub := domain.Subscription{
		ID:          uuid.New(),
		ServiceName: "Spotify",
		Price:       300,
		UserID:      uuid.New(),
		StartDate:   startDate,
		EndDate:     &endDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Act
	resp := toResponse(sub)

	// Assert

	if resp.StartDate != "01-2025" {
		t.Errorf("expected StartDate 01-2025, got %s", resp.StartDate)
	}

	if resp.EndDate != "03-2025" {
		t.Errorf("expected EndDate 03-2025, got %s", resp.EndDate)
	}
}

// ==============================================================
// ==============================================================
// parseSubscriptionRequest
// ==============================================================
// ==============================================================
func TestParseSubscriptionRequest_ValidWithoutEndDate(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
		EndDate:     "",
	}
	id := uuid.New()

	// Act
	sub, err := parseSubscriptionRequest(req, id)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sub.ID != id {
		t.Errorf("expected id %v, got %v", id, sub.ID)
	}
	if sub.ServiceName != "Netflix" {
		t.Errorf("expected service_name Netflix, got %s", sub.ServiceName)
	}
	if sub.Price != 500 {
		t.Errorf("expected price 500, got %d", sub.Price)
	}
	if sub.EndDate != nil {
		t.Errorf("expected end_date to be nil")
	}
}

func TestParseSubscriptionRequest_ValidWithEndDate(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Spotify",
		Price:       300,
		UserID:      uuid.New().String(),
		StartDate:   "01-2025",
		EndDate:     "03-2025",
	}
	id := uuid.New()

	// Act
	sub, err := parseSubscriptionRequest(req, id)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sub.EndDate == nil {
		t.Fatalf("expected end_date to be set")
	}

	expectedEnd := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	if !sub.EndDate.Equal(expectedEnd) {
		t.Errorf("expected end_date %v, got %v", expectedEnd, sub.EndDate)
	}
}

func TestParseSubscriptionRequest_EmptyServiceName(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "   ",
		Price:       100,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	// Act
	_, err := parseSubscriptionRequest(req, uuid.New())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseSubscriptionRequest_NegativePrice(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Netflix",
		Price:       -10,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
	}

	// Act
	_, err := parseSubscriptionRequest(req, uuid.New())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseSubscriptionRequest_InvalidUserID(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      "not-a-uuid",
		StartDate:   "07-2025",
	}

	// Act
	_, err := parseSubscriptionRequest(req, uuid.New())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseSubscriptionRequest_InvalidStartDate(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      uuid.New().String(),
		StartDate:   "2025-07",
	}

	// Act
	_, err := parseSubscriptionRequest(req, uuid.New())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseSubscriptionRequest_EndDateBeforeStartDate(t *testing.T) {
	// Arrange
	req := SubscriptionRequest{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      uuid.New().String(),
		StartDate:   "07-2025",
		EndDate:     "06-2025",
	}

	// Act
	_, err := parseSubscriptionRequest(req, uuid.New())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ==============================================================
// ==============================================================
// mapErrorToHTTP
// ==============================================================
// ==============================================================
func TestMapErrorToHTTP_Nil(t *testing.T) {
	// Act
	code, msg := mapErrorToHTTP(nil)

	// Assert
	if code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, code)
	}
	if msg != "" {
		t.Errorf("expected empty msg, got %q", msg)
	}
}

func TestMapErrorToHTTP_NotFound(t *testing.T) {
	// Arrange
	err := postgres.ErrNotFound

	// Act
	code, msg := mapErrorToHTTP(err)

	// Assert
	if code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, code)
	}
	if msg != "not found" {
		t.Errorf("expected %q, got %q", "not found", msg)
	}
}

func TestMapErrorToHTTP_NotFoundWrappedWithFmt(t *testing.T) {
	// Arrange
	wrapped := fmt.Errorf("query failed: %w", postgres.ErrNotFound)

	// Act
	code, msg := mapErrorToHTTP(wrapped)

	// Assert
	if code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, code)
	}
	if msg != "not found" {
		t.Errorf("expected %q, got %q", "not found", msg)
	}
}

func TestMapErrorToHTTP_DeadlineExceeded(t *testing.T) {
	// Arrange
	err := context.DeadlineExceeded

	// Act
	code, msg := mapErrorToHTTP(err)

	// Assert
	if code != http.StatusGatewayTimeout {
		t.Errorf("expected %d, got %d", http.StatusGatewayTimeout, code)
	}
	if msg != "timeout" {
		t.Errorf("expected %q, got %q", "timeout", msg)
	}
}

func TestMapErrorToHTTP_DeadlineExceededWrapped(t *testing.T) {
	// Arrange
	wrapped := errors.Join(context.DeadlineExceeded)

	// Act
	code, msg := mapErrorToHTTP(wrapped)

	// Assert
	if code != http.StatusGatewayTimeout {
		t.Errorf("expected %d, got %d", http.StatusGatewayTimeout, code)
	}
	if msg != "timeout" {
		t.Errorf("expected %q, got %q", "timeout", msg)
	}
}

func TestMapErrorToHTTP_DefaultInternal(t *testing.T) {
	// Arrange
	err := errors.New("some db error")

	// Act
	code, msg := mapErrorToHTTP(err)

	// Assert
	if code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, code)
	}
	if msg != "db error" {
		t.Errorf("expected %q, got %q", "db error", msg)
	}
}
