// Copyright 2016 Google Inc. All Rights Reserved.
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

package hilbert

import (
	"math/rand"
	"testing"
)

const peanoBenchmarkN64 = uint64(81)

// Test cases below assume N=9
var peanoTestCases64 = []struct {
	d, x, y uint64
}{
	{0, 0, 0},
	{1, 0, 1},
	{2, 0, 2},
	{3, 1, 2},
	{4, 1, 1},
	{5, 1, 0},
	{6, 2, 0},
	{7, 2, 1},
	{8, 2, 2},
	{9, 2, 3},
	// TODO Add more
}

func TestPeanoNewErrors64(t *testing.T) {
	var newTestCases = []struct {
		n    uint64
		want error
	}{
		{0, ErrNotPositive},
		{2, ErrNotPowerOfThree},
		{4, ErrNotPowerOfThree},
	}

	for _, tc := range newTestCases {
		s, err := NewPeano64(tc.n)
		if s != nil || err != tc.want {
			t.Errorf("NewPeano64(%d) = (%+v, %q) did not fail want (?, %q)", tc.n, s, err, tc.want)
		}
	}
}

func TestPeanoMapRangeErrors64(t *testing.T) {
	var mapRangeTestCases = []struct {
		d       uint64
		wantErr error
	}{
		{0, nil},
		{80, nil},
		{81, ErrOutOfRange},
	}

	s, err := NewPeano64(9)
	if err != nil {
		t.Fatalf("NewPeano64(9) failed: %s", err)
	}

	for _, tc := range mapRangeTestCases {
		if _, _, err = s.Map(tc.d); err != tc.wantErr {
			t.Errorf("Map(%d) = %q want %q", tc.d, tc.wantErr, err)
		}
	}
}

/*
func TestPeanoMapInverseRangeErrors64(t *testing.T) {
	var mapInverseRangeTestCases = []struct {
		x, y    uint64
		wantErr error
	}{
		{0, 0, nil},
		{15, 15, nil},
		{16, 0, ErrOutOfRange},
		{0, 16, ErrOutOfRange},
	}

	s, err := New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range mapInverseRangeTestCases {
		if _, err = s.MapInverse(tc.x, tc.y); err != tc.wantErr {
			t.Errorf("MapInverse(%d, %d) did not fail, want %q, got %q", tc.x, tc.y, tc.wantErr, err)
		}
	}
}
*/

func TestPeanoSmallMap64(t *testing.T) {
	s, err := NewPeano64(1)
	if err != nil {
		t.Fatalf("NewPeano(1) failed: %s", err)
	}

	x, y, err := s.Map(0)
	if err != nil {
		t.Errorf("Map(0) returned error: %s", err)
	}
	if x != 0 || y != 0 {
		t.Errorf("Map(0) = (%d, %d) want (0, 0)", x, y)
	}

	/*
		// TODO Test when MapInverse is implemented
		d, err := s.MapInverse(0, 0)
		if err != nil {
			t.Errorf("MapInverse(0,0) returned error: %s", err)
		}
		if d != 0 {
			t.Errorf("MapInverse(0, 0) failed, want 0, got %d", d)
		}
	*/
}

func TestPeanoMap64(t *testing.T) {
	s, err := NewPeano64(9)
	if err != nil {
		t.Fatalf("NewPeano(9) failed: %s", err)
	}

	for _, tc := range peanoTestCases64 {
		x, y, err := s.Map(tc.d)
		if err != nil {
			t.Errorf("Map(%d) returned error: %s", tc.d, err)
		}
		if x != tc.x || y != tc.y {
			t.Errorf("Map(%d) = (%d, %d) want (%d, %d)", tc.d, x, y, tc.x, tc.y)
		}
	}
}

/*
func TestPeanoMapInverse64(t *testing.T) {
	s, err := New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range testCases64 {
		d, err := s.MapInverse(tc.x, tc.y)
		if err != nil {
			t.Errorf("MapInverse(%d, %d) returned error: %s", tc.x, tc.y, err)
		}
		if d != tc.d {
			t.Errorf("MapInverse(%d, %d) failed, want %d, got %d", tc.x, tc.y, tc.d, d)
		}
	}
}

func TestPeanoAllMapValues64(t *testing.T) {
	s, err := New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for d := 0; d < s.N*s.N; d++ {
		// Map forwards and then back
		x, y, err := s.Map(d)
		if err != nil {
			t.Errorf("Map(%d) returned error: %s", d, err)
		}
		if x >= s.N || y >= s.N {
			t.Errorf("Map(%d) returned x,y out of range: (%d, %d)", d, x, y)
		}

		dPrime, err := s.MapInverse(x, y)
		if err != nil {
			t.Errorf("MapInverse(%d, %d) returned error: %s", x, y, err)
		}
		if d != dPrime {
			t.Errorf("Failed Map(%d) -> MapInverse(%d, %d) -> %d", d, x, y, dPrime)
		}
	}
}
*/
func BenchmarkPeanoMap64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewPeano64(peanoBenchmarkN)
		if err != nil {
			b.Fatalf("NewPeano64(%d) failed: %s", peanoBenchmarkN64, err)
		}
		for d := uint64(0); d < peanoBenchmarkN64*peanoBenchmarkN64; d++ {
			s.Map(d)
		}
	}
}

func BenchmarkPeanoMapRandom64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewPeano64(peanoBenchmarkN64)
		if err != nil {
			b.Fatalf("NewPeano64(%d) failed: %s", peanoBenchmarkN64, err)
		}
		for d := uint64(0); d < peanoBenchmarkN64*peanoBenchmarkN64; d++ {
			rd := rand.Intn(int(peanoBenchmarkN64 * peanoBenchmarkN64)) // Pick a random d
			s.Map(uint64(rd))
		}
	}
}

/*
func BenchmarkPeanoMapInverse64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := New(benchmarkN64)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}

		for x := 0; x < benchmarkN64; x++ {
			for y := 0; y < benchmarkN64; y++ {
				s.MapInverse(x, y)
			}
		}
	}
}
*/

func TestIsPow364(t *testing.T) {
	testCases := []struct {
		in   float64
		want bool
	}{
		{-1, false},
		{0, false},
		{1, true},
		{2, false},
		{3, true},
		{3.1, false},
		{4, false},
		{5, false},
		{8.9999, false},
		{9, true},
		{9.00001, false},
		{27, true},
		{59049, true},
	}

	for _, tc := range testCases {
		got := isPow3(tc.in)
		if got != tc.want {
			t.Errorf("isPow3(%f) = %t want %t", tc.in, got, tc.want)
		}
	}
}
