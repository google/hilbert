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
	"testing"
)

/*
import (
	"math/rand"
	"testing"

	"github.com/google/hilbert"
)
const benchmarkN = 81
*/

// Test cases below assume N=9
var peanoTestCases = []struct {
	d, x, y int
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
}

/*
func TestNewErrors(t *testing.T) {
	var newTestCases = []struct {
		n       int
		wantErr error
	}{
		{-1, hilbert.ErrNotPositive},
		{0, hilbert.ErrNotPositive},
		{3, hilbert.ErrNotPowerOfTwo},
		{5, hilbert.ErrNotPowerOfTwo},
	}

	for _, tc := range newTestCases {
		s, err := hilbert.New(tc.n)
		if s != nil || err != tc.wantErr {
			t.Errorf("New(%d) did not fail, want %q, got (%+v, %q)", tc.n, tc.wantErr, s, err)
		}
	}
}

func TestMapRangeErrors(t *testing.T) {
	var mapRangeTestCases = []struct {
		d       int
		wantErr error
	}{
		{-1, hilbert.ErrOutOfRange},
		{0, nil},
		{255, nil},
		{256, hilbert.ErrOutOfRange},
	}

	s, err := hilbert.New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range mapRangeTestCases {
		if _, _, err = s.Map(tc.d); err != tc.wantErr {
			t.Errorf("Map(%d) did not fail, want %q, got %q", tc.d, tc.wantErr, err)
		}
	}
}

func TestMapInverseRangeErrors(t *testing.T) {
	var mapInverseRangeTestCases = []struct {
		x, y    int
		wantErr error
	}{
		{0, 0, nil},
		{15, 15, nil},
		{-1, 0, hilbert.ErrOutOfRange},
		{0, -1, hilbert.ErrOutOfRange},
		{16, 0, hilbert.ErrOutOfRange},
		{0, 16, hilbert.ErrOutOfRange},
	}

	s, err := hilbert.New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range mapInverseRangeTestCases {
		if _, err = s.MapInverse(tc.x, tc.y); err != tc.wantErr {
			t.Errorf("MapInverse(%d, %d) did not fail, want %q, got %q", tc.x, tc.y, tc.wantErr, err)
		}
	}
}

func TestSmallMap(t *testing.T) {
	s, err := hilbert.New(1)
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
*/
func TestPeanoMap(t *testing.T) {
	s, err := NewPeano(9)
	if err != nil {
		t.Fatalf("Failed to create peano space: %s", err)
	}

	for _, tc := range peanoTestCases {
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
func TestMapInverse(t *testing.T) {
	s, err := hilbert.New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, tc := range testCases {
		d, err := s.MapInverse(tc.x, tc.y)
		if err != nil {
			t.Errorf("MapInverse(%d, %d) returned error: %s", tc.x, tc.y, err)
		}
		if d != tc.d {
			t.Errorf("MapInverse(%d, %d) failed, want %d, got %d", tc.x, tc.y, tc.d, d)
		}
	}
}

func TestAllMapValues(t *testing.T) {
	s, err := hilbert.New(16)
	if err != nil {
		t.Fatalf("Failed to create hibert space: %s", err)
	}

	for d := 0; d < s.N*s.N; d++ {
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

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}
		for d := 0; d < benchmarkN*benchmarkN; d++ {
			s.Map(d)
		}
	}
}

func BenchmarkMapRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}
		for d := 0; d < benchmarkN*benchmarkN; d++ {
			rd := rand.Intn(benchmarkN * benchmarkN) // Pick a random d
			s.Map(rd)
		}
	}
}

func BenchmarkMapInverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			b.Fatalf("Failed to create hibert space: %s", err)
		}

		for x := 0; x < benchmarkN; x++ {
			for y := 0; y < benchmarkN; y++ {
				s.MapInverse(x, y)
			}
		}
	}
}
*/

func TestIsPow3(t *testing.T) {
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
