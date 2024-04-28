package service

import (
	"testing"
)

func TestValidateDonation(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		amount := 200000.0

		got := ValidateDonation(amount)

		if got != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, got)
		}
	})
	t.Run("should return 10000", func(t *testing.T) {
		amount := 10000.0

		got := ValidateDonation(amount)

		if got != 10000.0 {
			t.Errorf("expected %f, but got %f", 10000.0, got)
		}
	})
}
