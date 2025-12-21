package handlers

import (
	"testing"
	"time"
)

// ====================================
// monthsInclusive
// ====================================

// TestMonthsInclusive_SameMonth verifies that the same month counts as one month.
func TestMonthsInclusive_SameMonth(t *testing.T) {
	// Arrange
	start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := monthsInclusive(start, end)

	// Assert
	if result != 1 {
		t.Errorf("expected 1 month, got %d", result)
	}
}

// TestMonthsInclusive_SameYear verifies counting across several months in the same year.
func TestMonthsInclusive_SameYear(t *testing.T) {
	// Arrange
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := monthsInclusive(start, end)

	// Assert
	if result != 3 {
		t.Errorf("expected 3 months, got %d", result)
	}
}

// TestMonthsInclusive_CrossYear verifies month counting across year boundary.
func TestMonthsInclusive_CrossYear(t *testing.T) {
	// Arrange
	start := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := monthsInclusive(start, end)

	// Assert
	if result != 4 {
		t.Errorf("expected 4 months, got %d", result)
	}
}

// ====================================
// maxTime
// ====================================

// TestMaxTime_After checks that the later time is returned.
func TestMaxTime_After(t *testing.T) {
	// Arrange
	a := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := maxTime(a, b)

	// Assert
	if !result.Equal(a) {
		t.Errorf("expected %v, got %v", a, result)
	}
}

// TestMaxTime_Before checks that the later time is returned
func TestMaxTime_Before(t *testing.T) {
	// Arrange
	a := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := maxTime(a, b)

	// Assert
	if !result.Equal(b) {
		t.Errorf("expected %v, got %v", b, result)
	}
}

// TestMaxTime_Equal checks behavior when times are equal.
func TestMaxTime_Equal(t *testing.T) {
	// Arrange
	a := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := maxTime(a, b)

	// Assert
	if !result.Equal(b) {
		t.Errorf("expected %v, got %v", b, result)
	}
}

// ====================================
// minTime
// ====================================

// TestMinTime_Before checks that the earlier time is returned.
func TestMinTime_Before(t *testing.T) {
	// Arrange
	a := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := minTime(a, b)

	// Assert
	if !result.Equal(a) {
		t.Errorf("expected %v, got %v", a, result)
	}
}

// TestMinTime_After checks that the earlier time is returned.
func TestMinTime_After(t *testing.T) {
	// Arrange
	a := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := minTime(a, b)

	// Assert
	if !result.Equal(b) {
		t.Errorf("expected %v, got %v", b, result)
	}
}

// TestMinTime_Equal checks behavior when times are equal.
func TestMinTime_Equal(t *testing.T) {
	// Arrange
	a := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	b := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := minTime(a, b)

	// Assert
	if !result.Equal(b) {
		t.Errorf("expected %v, got %v", b, result)
	}
}
