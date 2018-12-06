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

package hilbert

import (
	"math/rand"
	"testing"
)

const benchmarkN64 = uint64(32)

// Test cases below assume N=16
var testCases64 = []struct {
	d, x, y uint64
}{
	{0, 0, 0},
	{16, 4, 0},
	{32, 4, 4},
	{48, 3, 7},
	{64, 0, 8},
	{80, 0, 12},
	{96, 4, 12},
	{112, 7, 11},
	{128, 8, 8},
	{144, 8, 12},
	{160, 12, 12},
	{170, 15, 15},
	{176, 15, 11},
	{192, 15, 7},
	{208, 11, 7},
	{224, 11, 3},
	{240, 12, 0},
	{255, 15, 0},
}

func TestNewErrors64(t *testing.T) {
	var newTestCases = []struct {
		n       uint64
		wantErr error
	}{
		{0, ErrNotPositive},
		{3, ErrNotPowerOfTwo},
		{5, ErrNotPowerOfTwo},
	}

	for _, tc := range newTestCases {
		s, err := NewHilbert64(tc.n)
		if s != nil || err != tc.wantErr {
			t.Errorf("NewHilbert64(%d) did not fail, want %q, got (%+v, %q)", tc.n, tc.wantErr, s, err)
		}
	}
}

func TestMapRangeErrors64(t *testing.T) {
	var mapRangeTestCases = []struct {
		d       uint64
		wantErr error
	}{
		{0, nil},
		{255, nil},
		{256, ErrOutOfRange},
	}

	s, err := NewHilbert64(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range mapRangeTestCases {
		if _, _, err = s.Map(tc.d); err != tc.wantErr {
			t.Errorf("Map(%d) did not fail, want %q, got %q", tc.d, tc.wantErr, err)
		}
	}
}

func TestMapInverseRangeErrors64(t *testing.T) {
	var mapInverseRangeTestCases = []struct {
		x, y    uint64
		wantErr error
	}{
		{0, 0, nil},
		{15, 15, nil},
		{16, 0, ErrOutOfRange},
		{0, 16, ErrOutOfRange},
	}

	s, err := NewHilbert64(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range mapInverseRangeTestCases {
		if _, err = s.MapInverse(tc.x, tc.y); err != tc.wantErr {
			t.Errorf("MapInverse(%d, %d) did not fail, want %q, got %q", tc.x, tc.y, tc.wantErr, err)
		}
	}
}

func TestSmallMap64(t *testing.T) {
	s, err := NewHilbert64(1)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	x, y, err := s.Map(0)
	if err != nil {
		t.Errorf("Map(0) returned error: %s", err)
	}
	if x != 0 || y != 0 {
		t.Errorf("Map(0) failed, want (0, 0), got (%d, %d)", x, y)
	}

	d, err := s.MapInverse(0, 0)
	if err != nil {
		t.Errorf("MapInverse(0,0) returned error: %s", err)
	}
	if d != 0 {
		t.Errorf("MapInverse(0, 0) failed, want 0, got %d", d)
	}
}

func TestMap64(t *testing.T) {
	s, err := NewHilbert64(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range testCases64 {
		x, y, err := s.Map(tc.d)
		if err != nil {
			t.Errorf("Map(%d) returned error: %s", tc.d, err)
		}
		if x != tc.x || y != tc.y {
			t.Errorf("Map(%d) failed, want (%d, %d), got (%d, %d)", tc.d, tc.x, tc.y, x, y)
		}
	}
}

func TestMapInverse64(t *testing.T) {
	s, err := NewHilbert64(16)
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

func TestAllMapValues64(t *testing.T) {
	s, err := NewHilbert64(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for d := uint64(0); d < s.N*s.N; d++ {
		// Map forwards and then back
		x, y, err := s.Map(d)
		if err != nil {
			t.Errorf("Map(%d) returned error: %s", d, err)
		}
		if x < 0 || x >= s.N || y < 0 || y >= s.N {
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

func BenchmarkMap64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewHilbert64(benchmarkN64)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}
		for d := uint64(0); d < benchmarkN64*benchmarkN64; d++ {
			s.Map(d)
		}
	}
}

func BenchmarkMapRandom64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewHilbert64(benchmarkN64)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}
		for d := uint64(0); d < benchmarkN64*benchmarkN64; d++ {
			// Note that we're going to coerce this to an implied 32-bit type 
			// as required by Intn() but that the value is tiny and will always 
			// be within the range.
			rd := rand.Intn(int(benchmarkN64 * benchmarkN64)) // Pick a random d
			s.Map(uint64(rd))
		}
	}
}

func BenchmarkMapInverse64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewHilbert64(benchmarkN64)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}

		for x := uint64(0); x < benchmarkN64; x++ {
			for y := uint64(0); y < benchmarkN64; y++ {
				s.MapInverse(x, y)
			}
		}
	}
}
