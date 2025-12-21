package utils

import (
	"testing"
	"time"
)

// ============================
// ParseMonthYear
// ============================

func TestParseMonthYear_ValidJanuary(t *testing.T) {
	// Arrange
	input := "01-2025"

	// Act
	result, err := ParseMonthYear(input)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseMonthYear_ValidDecember(t *testing.T) {
	// Arrange
	input := "12-2024"

	// Act
	result, err := ParseMonthYear(input)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseMonthYear_ValidJuly(t *testing.T) {
	// Arrange
	input := "07-2025"

	// Act
	result, err := ParseMonthYear(input)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseMonthYear_InvalidFormat(t *testing.T) {
	// Arrange
	input := "2025-07"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_InvalidMonthHigh(t *testing.T) {
	// Arrange
	input := "13-2025"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_InvalidMonthZero(t *testing.T) {
	// Arrange
	input := "00-2025"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_InvalidMonthNotNumber(t *testing.T) {
	// Arrange
	input := "aa-2025"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_InvalidYearNotNumber(t *testing.T) {
	// Arrange
	input := "07-aa"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_YearTooSmall(t *testing.T) {
	// Arrange
	input := "07-1899"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMonthYear_YearTooLarge(t *testing.T) {
	// Arrange
	input := "07-2500"

	// Act
	_, err := ParseMonthYear(input)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ===============================
// FormatMonthYear
// ===============================

func TestFormatMonthYear_StandardDate(t *testing.T) {
	// Arrange
	input := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	// Act
	result := FormatMonthYear(input)

	// Assert
	if result != "07-2025" {
		t.Errorf("expected 07-2025, got %s", result)
	}
}

func TestFormatMonthYear_NormalizesDay(t *testing.T) {
	// Arrange
	input := time.Date(2025, 7, 15, 10, 30, 0, 0, time.UTC)

	// Act
	result := FormatMonthYear(input)

	// Assert
	if result != "07-2025" {
		t.Errorf("expected 07-2025, got %s", result)
	}
}

func TestFormatMonthYear_LeadingZero(t *testing.T) {
	// Arrange
	input := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	// Act
	result := FormatMonthYear(input)

	// Assert
	if result != "01-2024" {
		t.Errorf("expected 01-2024, got %s", result)
	}
}

func TestFormatMonthYear_EndOfYear(t *testing.T) {
	// Arrange
	input := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Act
	result := FormatMonthYear(input)

	// Assert
	if result != "12-2024" {
		t.Errorf("expected 12-2024, got %s", result)
	}
}
