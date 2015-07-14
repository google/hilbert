package hilbert_test

func Example() {
	// Create a Hilbert Curve for mapping to and from a 16 by 16 space
	s, err := New(16)

	// Now map one dimension numbers in the range [0, N*N-1], to an x,y
	// coordinate on the curve where both x and y are in the range [0, N-1].
	x, y, err := s.Map(96)

	// Also map back from (x,y) to t
	t, err := s.MapInverse(x, y)
	// Output:
	// x = 4, y = 12
	// t = 96
}
