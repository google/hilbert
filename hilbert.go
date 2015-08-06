// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package hilbert provides mapping of values to and from Hilbert curves.
//
// Converted from the code available on Wikipedia code, with additional help from:
//  * https://en.wikipedia.org/wiki/Hilbert_curve
//  * http://bit-player.org/2013/mapping-the-hilbert-curve
//
package hilbert

import (
	"errors"
)

// Errors returned when validating input.
var (
	ErrLessThanZero  = errors.New("N must be greater than zero")
	ErrNotPowerOfTwo = errors.New("N must be a power of two")
	ErrOutOfRange    = errors.New("Value is out of range")
)

// Space represents a 2D Hilbert space of order N for mapping to and from.
type Space struct {
	N int
}

// New returns a Hilbert space which maps integers to and from the curve.
// n must be a power of two.
func New(n int) (*Space, error) {
	if n <= 0 {
		return nil, ErrLessThanZero
	}

	// Test if power of two
	if (n & (n - 1)) != 0 {
		return nil, ErrNotPowerOfTwo
	}

	return &Space{
		N: n,
	}, nil
}

func i2b(i int) bool {
	return i != 0
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Map transforms a dimension value, t, in the range [0, n^2-1] to coordinates on the Hilbert
// curve in the two-dimension space, where x and y are within [0,n-1].
func (s *Space) Map(t int) (x, y int, err error) {
	if t < 0 || t >= s.N*s.N {
		return -1, -1, ErrOutOfRange
	}

	x = 0
	y = 0

	for i := 1; i < s.N; i = i * 2 {
		rx := i2b(1 & (t / 2)) // TODO make more go'ish
		ry := i2b(1 & (t ^ b2i(rx)))
		x, y = rot(i, x, y, rx, ry)

		x = x + i*b2i(rx)
		y = y + i*b2i(ry)
		t /= 4
	}

	return
}

// MapInverse transform coordinates on Hilbert Curve from (x,y) to t.
func (s *Space) MapInverse(x, y int) (t int, err error) {
	if x < 0 || x >= s.N || y < 0 || y >= s.N {
		return -1, ErrOutOfRange
	}

	t = 0
	for i := s.N / 2; i > 0; i = i / 2 {
		rx := (x & i) > 0
		ry := (y & i) > 0
		t += i * i * ((3 * b2i(rx)) ^ b2i(ry))
		x, y = rot(i, x, y, rx, ry)
	}

	return
}

// Rotate/flip a quadrant appropriately
func rot(n, x, y int, rx, ry bool) (int, int) {
	if !ry {
		if rx {
			x = n - 1 - x
			y = n - 1 - y
		}

		//Swap x and y
		x, y = y, x
	}
	return x, y
}
