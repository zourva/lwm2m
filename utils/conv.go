package utils

import "fmt"

// IntToStr int/uint/int32/uint32/int64/uint64
// are expected.
func IntToStr[T int | uint | int32 | uint32 | int64 | uint64](t T) string {
	return fmt.Sprintf("%d", t)
}

// FloatToStr float32 and float64 are expected.
func FloatToStr[T float32 | float64](t T, precision int) string {
	if precision < 0 {
		precision = 6
	}

	ff := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(ff, t)
}
