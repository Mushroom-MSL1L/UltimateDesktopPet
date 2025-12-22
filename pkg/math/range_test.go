package math

import "testing"

func TestInRange(t *testing.T) {
	tests := []struct {
		name     string
		x, max   int16
		min      int16
		expected int16
	}{
		{name: "below_min", x: -1, max: 100, min: 0, expected: 0},
		{name: "above_max", x: 101, max: 100, min: 0, expected: 100},
		{name: "within_range", x: 50, max: 100, min: 0, expected: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InRange(tt.x, tt.max, tt.min); got != tt.expected {
				t.Fatalf("InRange(%d, %d, %d) = %d, want %d", tt.x, tt.max, tt.min, got, tt.expected)
			}
		})
	}
}
