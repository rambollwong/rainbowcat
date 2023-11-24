package util

import (
	"testing"
)

func TestParseToBytesSize(t *testing.T) {
	tests := []struct {
		input    string
		base     int64
		expected int64
	}{
		{"100", 1024, 100},
		{"1K", 1024, 1024},
		{"1M", 1024, 1048576},
		{"1G", 1024, 1073741824},
		{"1T", 1024, 1099511627776},
		{"500B", 1024, 500},
		{"2.5K", 1024, 2560},

		{"100", 1000, 100},
		{"1K", 1000, 1000},
		{"1M", 1000, 1000000},
		{"1G", 1000, 1000000000},
		{"1T", 1000, 1000000000000},
		{"500B", 1000, 500},
		{"2.5K", 1000, 2500},
	}

	for _, test := range tests {
		result, err := ParseToBytesSize(test.input, test.base)
		if err != nil {
			t.Errorf("Error parsing size string '%s': %s", test.input, err)
		}

		if result != test.expected {
			t.Errorf("Size mismatch for input '%s'. Expected: %d, Got: %d", test.input, test.expected, result)
		}
	}

	_, err := ParseToBytesSize("2.5.5", 1000)
	if err == nil {
		t.Errorf("Should return an error")
	}
}
