package hilbert

import "errors"

// Errors returned when validating input.
var (
	ErrNotPositive     = errors.New("N must be greater than zero")
	ErrNotPowerOfTwo   = errors.New("N must be a power of two")
	ErrNotPowerOfThree = errors.New("N must be a power of three")
	ErrOutOfRange      = errors.New("value is out of range")
)

// SpaceFilling represents a space-filling curve that can map points from one dimensions to two.
type SpaceFilling interface {
	// Map transforms a one dimension value, t, in the range [0, n^2-1] to coordinates on the
	// curve in the two-dimension space, where x and y are within [0,n-1].
	Map(t int) (x, y int, err error)

	// MapInverse transform coordinates on the curve from (x,y) to t.
	MapInverse(x, y int) (t int, err error)

	// GetDimensions returns the width and height of the 2D space.
	GetDimensions() (x, y int)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
