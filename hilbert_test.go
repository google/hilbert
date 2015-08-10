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

package hilbert_test

import (
	"testing"

	"github.com/google/hilbert"
	"math/rand"
)

const benchmarkN = 32

// Test cases when N=16
var testCases = []struct {
	t, x, y int
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

var newTestCases = []struct {
	n           int
	expectedErr error
}{
	{-1, hilbert.ErrLessThanZero},
	{0, hilbert.ErrLessThanZero},
	{3, hilbert.ErrNotPowerOfTwo},
	{5, hilbert.ErrNotPowerOfTwo},
}

func TestNewErrors(test *testing.T) {

	for _, testCase := range newTestCases {
		s, err := hilbert.New(testCase.n)
		if s != nil || err != testCase.expectedErr {
			test.Errorf(
				"New(%d) did not fail, expected '%s', got (%+v, %v)",
				testCase.n, testCase.expectedErr, s, err)
		}
	}
}

func TestSmallMap(test *testing.T) {
	s, err := hilbert.New(1)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	x, y, err := s.Map(0)
	if x != 0 || y != 0 {
		test.Errorf(
			"Map(%d) failed, expected (%d, %d), got (%d, %d)",
			0, 0, 0, x, y)
	}

	t, err := s.MapInverse(0, 0)
	if t != 0 {
		test.Errorf(
			"MapInverse(%d, %d) failed, expected %d, got %d",
			0, 0, 0, t)
	}
}

func TestMap(test *testing.T) {
	s, err := hilbert.New(16)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, testCase := range testCases {
		x, y, err := s.Map(testCase.t)
		if err != nil {
			test.Errorf("Map(%d) returned error: %s", testCase.t, err.Error())

		} else if x != testCase.x || y != testCase.y {
			test.Errorf(
				"Map(%d) failed, expected (%d, %d), got (%d, %d)",
				testCase.t, testCase.x, testCase.y, x, y)
		}
	}
}

func TestMapInverse(test *testing.T) {
	s, err := hilbert.New(16)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, testCase := range testCases {
		t, err := s.MapInverse(testCase.x, testCase.y)
		if err != nil {
			test.Errorf("MapInverse(%d, %d) returned error: %s", testCase.x, testCase.y, err.Error())
		}
		if t != testCase.t {
			test.Errorf(
				"MapInverse(%d, %d) failed, expected %d, got %d",
				testCase.x, testCase.y, testCase.t, t)
		}
	}
}

func TestAllMapValues(test *testing.T) {
	s, err := hilbert.New(16)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	for t := 0; t < s.N*s.N; t++ {
		// Map forwards and then back
		x, y, err := s.Map(t)
		if err != nil {
			test.Errorf("Map(%d) returned error: %s", t, err.Error())
		}

		if x < 0 || x >= s.N || y < 0 || y >= s.N {
			test.Errorf(
				"Map(%d) returned x,y out of range: (%d, %d)",
				t, x, y)
		}

		tPrime, err := s.MapInverse(x, y)
		if err != nil {
			test.Errorf("MapInverse(%d, %d) returned error: %s", x, y, err.Error())
		}
		if t != tPrime {
			test.Errorf(
				"Failed Map(%d) -> MapInverse(%d, %d) -> %d",
				t, x, y, tPrime)
		}

	}
}

func BenchmarkMap(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			benchmark.Fatalf("Failed to create hibert space: %s", err)
		}
		for t := 0; t < benchmarkN*benchmarkN; t++ {
			s.Map(t)
		}
	}
}

func BenchmarkMapRandom(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			benchmark.Fatalf("Failed to create hibert space: %s", err)
		}
		for t := 0; t < benchmarkN*benchmarkN; t++ {
			rt := rand.Intn(benchmarkN * benchmarkN) // Pick a random t
			s.Map(rt)
		}
	}
}

func BenchmarkMapInverse(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		s, err := hilbert.New(benchmarkN)
		if err != nil {
			benchmark.Fatalf("Failed to create hibert space: %s", err)
		}

		for x := 0; x < benchmarkN; x++ {
			for y := 0; y < benchmarkN; y++ {
				s.MapInverse(x, y)
			}
		}
	}
}
