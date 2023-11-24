package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ParseToBytesSize converts a string with size and unit suffix to byte size.
//
// params:
//   - sizeStr: Size string, e.g., "1K", "500B", "2M", etc.
//   - base: Base used to calculate the multiplication factor for units, e.g., 1024.
//
// Returns the converted byte size and possible error.
func ParseToBytesSize(sizeStr string, base int64) (int64, error) {
	pattern := `^(\d+(?:\.\d+)?)([bBkKmMgGtT]?)$`
	regex := regexp.MustCompile(pattern)

	if !regex.Match([]byte(sizeStr)) {
		return 0, errors.New("invalid size string")
	}
	match := regex.FindStringSubmatch(sizeStr)
	if match == nil {
		return 0, errors.New("invalid size string")
	}

	valueStr := match[1]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}
	if len(match) > 2 {
		unit := strings.ToLower(match[2])
		switch unit {
		case "b":
			return int64(value), nil
		case "k":
			return int64(value * float64(base)), nil
		case "m":
			return int64(value * float64(base*base)), nil
		case "g":
			return int64(value * float64(base*base*base)), nil
		case "t":
			return int64(value * float64(base*base*base*base)), nil
		}
	}
	return int64(value), nil
}
