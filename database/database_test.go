package database

import (
	"testing"
)

func TestValidatePersonal(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		amount := 200000

		got := ValidatePersonal(float64(amount))

		if got != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, got)
		}
	})
	t.Run("should return 15000", func(t *testing.T) {
		amount := 15000

		got := ValidatePersonal(float64(amount))

		if got != 15000.0 {
			t.Errorf("expected %f, but got %f", 15000.0, got)
		}
	})
	t.Run("should return 10001", func(t *testing.T) {
		amount := 200

		got := ValidatePersonal(float64(amount))

		if got != 10001.0 {
			t.Errorf("expected %f, but got %f", 10001.0, got)
		}
	})
}

func TestValidateMaxKReceipt(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		amount := 200000

		got := ValidateMaxKReceipt(float64(amount))

		if got != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, got)
		}
	})
	t.Run("should return 1", func(t *testing.T) {
		amount := -15000

		got := ValidateMaxKReceipt(float64(amount))

		if got != 1.0 {
			t.Errorf("expected %f, but got %f", 1.0, got)
		}
	})
	t.Run("should return 1500", func(t *testing.T) {
		amount := 1500

		got := ValidateMaxKReceipt(float64(amount))

		if got != 1500.0 {
			t.Errorf("expected %f, but got %f", 1500.0, got)
		}
	})
}
