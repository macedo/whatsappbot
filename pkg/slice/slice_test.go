package slice_test

import (
	"testing"

	"github.com/macedo/whatsappbot/pkg/slice"
)

func TestContainsString(t *testing.T) {
	name := []string{"john", "doe"}

	if !slice.Contains[string](name, "john") {
		t.Errorf("expected return true, got false")
	}

	if slice.Contains[string](name, "macedo") {
		t.Errorf("expected return false, got true")
	}
}

func TestContainsInt(t *testing.T) {
	ids := []int{1, 2, 3}

	if !slice.Contains[int](ids, 1) {
		t.Errorf("expected return true, got false")
	}

	if slice.Contains[int](ids, 10) {
		t.Errorf("expected return false, got true")
	}
}
