package hilbert

// Peano represents a 2D Peano curve of order N for mapping to and from.
// Implements SpaceFilling interface.
type Peano64 struct {
	N uint64 // Always a power of three, and is the width/height of the space.
}

// NewPeano returns a new Peano space filling curve which maps integers to and from the curve.
// n must be a power of three.
func NewPeano64(n uint64) (*Peano64, error) {
	if n == 0 {
		return nil, ErrNotPositive
	}

	if !isPow3(float64(n)) {
		return nil, ErrNotPowerOfThree
	}

	return &Peano64{
		N: n,
	}, nil
}

// GetDimensions returns the width and height of the 2D space.
func (p *Peano64) GetDimensions() (uint64, uint64) {
	return p.N, p.N
}

// Map transforms a one dimension value, t, in the range [0, n^3-1] to coordinates on the Peano
// curve in the two-dimension space, where x and y are within [0,n-1].
func (p *Peano64) Map(t uint64) (x, y uint64, err error) {
	if t >= p.N*p.N {
		return 0, 0, ErrOutOfRange
	}

	for i := uint64(1); i < p.N; i = i * 3 {
		s := t % 9

		// rx/ry are the coordinates in the 3x3 grid
		rx := uint64(s / 3)
		ry := uint64(s % 3)
		if rx == 1 {
			ry = 2 - ry
		}

		// now based on depth rotate our points
		if i > 1 {
			x, y = p.rotate(i, x, y, s)
		}

		x += rx * i
		y += ry * i

		t /= 9
	}

	return x, y, nil
}

// rotate rotates the x and y coordinates depending on the current n depth.
func (p *Peano64) rotate(n, x, y, s uint64) (uint64, uint64) {

	if n == 1 {
		// Special case
		return x, y
	}

	n = n - 1
	switch s {
	case 0:
		return x, y // normal
	case 1:
		return n - x, y // fliph
	case 2:
		return x, y // normal
	case 3:
		return x, n - y // flipv
	case 4:
		return n - x, n - y // flipv and fliph
	case 5:
		return x, n - y // flipv
	case 6:
		return x, y // normal
	case 7:
		return n - x, y // fliph
	case 8:
		return x, y // normal
	}

	panic("assertion failure: this line should never be reached")
}

// MapInverse transform coordinates on the Peano curve from (x,y) to t.
// NOT IMPLEMENTED YET
func (p *Peano64) MapInverse(x, y uint64) (t uint64, err error) {
	if x >= p.N || y >= p.N {
		return 0, ErrOutOfRange
	}

	panic("Not finished")
	return 0, nil
}
