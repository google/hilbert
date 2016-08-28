package hilbert

// Peano represents a 2D Peano curve of order N for mapping to and from.
// Implements SpaceFilling interface.
type Peano struct {
	N int // Always a power of three, and is the width/height of the space.
}

// isPow3 returns true if n is a power of 3.
func isPow3(n float64) bool {
	// I wanted to do the following, but due to subtle floating point issues it didn't work
	// const ln3 = 1.098612288668109691395245236922525704647490557822749451734694333637494 // https://oeis.org/A002391
	//return n == math.Pow(3, math.Trunc(math.Log(n) / ln3))
	for n >= 1 {
		if n == 1 {
			return true
		}
		n = n / 3
	}
	return false
}

// NewPeano returns a new Peano space filling curve which maps integers to and from the curve.
// n must be a power of three.
func NewPeano(n int) (*Peano, error) {
	if n <= 0 {
		return nil, ErrNotPositive
	}

	if !isPow3(float64(n)) {
		return nil, ErrNotPowerOfThree
	}

	return &Peano{
		N: n,
	}, nil
}

// GetDimensions returns the width and height of the 2D space.
func (p *Peano) GetDimensions() (int, int) {
	return p.N, p.N
}

// Map transforms a one dimension value, t, in the range [0, n^3-1] to coordinates on the Peano
// curve in the two-dimension space, where x and y are within [0,n-1].
func (p *Peano) Map(t int) (x, y int, err error) {
	if t < 0 || t >= p.N*p.N {
		return -1, -1, ErrOutOfRange
	}

	for i := 1; i < p.N; i = i * 3 {
		s := t % 9

		// rx/ry are the coordinates in the 3x3 grid
		rx := int(s / 3)
		ry := int(s % 3)
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
func (p *Peano) rotate(n, x, y, s int) (int, int) {

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
func (p *Peano) MapInverse(x, y int) (t int, err error) {
	if x < 0 || x >= p.N || y < 0 || y >= p.N {
		return -1, ErrOutOfRange
	}

	panic("Not finished")
	return -1, nil
}
