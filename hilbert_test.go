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
	"math/rand"
	"testing"

	"github.com/google/hilbert"
)

const benchmarkN = 32

// Test cases below assume N=16
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

func TestNewErrors(test *testing.T) {
	var newTestCases = []struct {
		n           int
		expectedErr error
	}{
		{-1, hilbert.ErrNotPositive},
		{0, hilbert.ErrNotPositive},
		{3, hilbert.ErrNotPowerOfTwo},
		{5, hilbert.ErrNotPowerOfTwo},
	}

	for _, testCase := range newTestCases {
		s, err := hilbert.New(testCase.n)
		if s != nil || err != testCase.expectedErr {
			test.Errorf(
				"New(%d) did not fail, expected %q, got (%+v, %q)",
				testCase.n, testCase.expectedErr, s, err)
		}
	}
}

func TestMapRangeErrors(test *testing.T) {
	var mapRangeTestCases = []struct {
		t           int
		expectedErr error
	}{
		{0, nil},
		{-1, hilbert.ErrOutOfRange},
		{256, hilbert.ErrOutOfRange},
	}

	s, err := hilbert.New(16)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, testCase := range mapRangeTestCases {
		if _, _, err = s.Map(testCase.t); err != testCase.expectedErr {
			test.Errorf(
				"Map(%d) did not fail, expected %q, got %q",
				testCase.t, testCase.expectedErr, err)
		}
	}
}

func TestMapInverseRangeErrors(test *testing.T) {
	var mapInverseRangeTestCases = []struct {
		x, y        int
		expectedErr error
	}{
		{0, 0, nil},
		{-1, 0, hilbert.ErrOutOfRange},
		{0, -1, hilbert.ErrOutOfRange},
		{16, 0, hilbert.ErrOutOfRange},
		{0, 16, hilbert.ErrOutOfRange},
	}

	s, err := hilbert.New(16)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	for _, testCase := range mapInverseRangeTestCases {
		if _, err = s.MapInverse(testCase.x, testCase.y); err != testCase.expectedErr {
			test.Errorf(
				"MapInverse(%d, %d) did not fail, expected %q, got %q",
				testCase.x, testCase.y, testCase.expectedErr, err)
		}
	}
}

func TestSmallMap(test *testing.T) {
	s, err := hilbert.New(1)
	if err != nil {
		test.Fatalf("Failed to create hibert space: %s", err)
	}

	x, y, err := s.Map(0)
	if err != nil {
		test.Errorf("Map(0) returned error: %s", err)
	}
	if x != 0 || y != 0 {
		test.Errorf("Map(0) failed, expected (0, 0), got (%d, %d)", x, y)
	}

	t, err := s.MapInverse(0, 0)
	if err != nil {
		test.Errorf("MapInverse(0,0) returned error: %s", err)
	}
	if t != 0 {
		test.Errorf("MapInverse(0, 0) failed, expected 0, got %d", t)
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
			test.Errorf("Map(%d) returned error: %s", testCase.t, err)
		}
		if x != testCase.x || y != testCase.y {
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
			test.Errorf("MapInverse(%d, %d) returned error: %s", testCase.x, testCase.y, err)
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
			test.Errorf("Map(%d) returned error: %s", t, err)
		}
		if x < 0 || x >= s.N || y < 0 || y >= s.N {
			test.Errorf(
				"Map(%d) returned x,y out of range: (%d, %d)",
				t, x, y)
		}

		tPrime, err := s.MapInverse(x, y)
		if err != nil {
			test.Errorf("MapInverse(%d, %d) returned error: %s", x, y, err)
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
