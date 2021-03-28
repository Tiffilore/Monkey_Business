package minimalexample

import (
	"testing"
)

func access(i int) string {
	return []string{"a", "b"}[i]
}

func TestAccess(t *testing.T) {
	tests := []struct {
		testdatum int
		expected  string
	}{
		{0, "a"},
		{2, ""},
		{1, "c"},
		{4, ""},
	}

	for _, tt := range tests {
		testAccess(tt, t)
	}
}

func testAccess(tt struct {
	testdatum int
	expected  string
}, t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Runtime error for testdatum %v: %q", tt.testdatum, err)
		}
	}()
	result := access(tt.testdatum)
	if result != tt.expected {
		t.Errorf("%v evaluates to %q, though it should be %q", tt.testdatum, result, tt.expected)
	}

}
