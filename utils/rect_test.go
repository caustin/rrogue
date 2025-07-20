package utils

import "testing"

func TestNewRect(t *testing.T) {
	rect := NewRect(10, 20, 5, 3)

	if rect.X1 != 10 || rect.Y1 != 20 || rect.X2 != 15 || rect.Y2 != 23 {
		t.Errorf("NewRect(10, 20, 5, 3) = {%d, %d, %d, %d}, expected {10, 20, 15, 23}",
			rect.X1, rect.Y1, rect.X2, rect.Y2)
	}
}

func TestRectCenter(t *testing.T) {
	rect := NewRect(0, 0, 10, 10)
	x, y := rect.Center()

	if x != 5 || y != 5 {
		t.Errorf("Center() = (%d, %d), expected (5, 5)", x, y)
	}
}

func TestRectIntersect(t *testing.T) {
	tests := []struct {
		name     string
		rect1    Rect
		rect2    Rect
		expected bool
	}{
		{
			name:     "overlapping rectangles",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(2, 2, 5, 5),
			expected: true,
		},
		{
			name:     "adjacent rectangles (touching edges)",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(5, 0, 5, 5),
			expected: true,
		},
		{
			name:     "completely separate rectangles",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(10, 10, 5, 5),
			expected: false,
		},
		{
			name:     "one rectangle inside another",
			rect1:    NewRect(0, 0, 10, 10),
			rect2:    NewRect(2, 2, 3, 3),
			expected: true,
		},
		{
			name:     "rectangles separated vertically",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(2, 6, 5, 5),
			expected: false,
		},
		{
			name:     "rectangles separated horizontally",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(6, 2, 5, 5),
			expected: false,
		},
		{
			name:     "identical rectangles",
			rect1:    NewRect(5, 5, 3, 3),
			rect2:    NewRect(5, 5, 3, 3),
			expected: true,
		},
		{
			name:     "corner touching",
			rect1:    NewRect(0, 0, 5, 5),
			rect2:    NewRect(5, 5, 5, 5),
			expected: true,
		},
		{
			name:     "diagonal separation",
			rect1:    NewRect(0, 0, 3, 3),
			rect2:    NewRect(4, 4, 3, 3),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Intersect(tt.rect2)
			if result != tt.expected {
				t.Errorf("Intersect() = %v, expected %v for %s", result, tt.expected, tt.name)
			}

			// Test symmetry - intersection should work both ways
			result2 := tt.rect2.Intersect(tt.rect1)
			if result2 != tt.expected {
				t.Errorf("Intersect() symmetry failed: %v != %v for %s", result, result2, tt.name)
			}
		})
	}
}
