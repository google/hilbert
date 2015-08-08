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
	"fmt"

	"github.com/google/hilbert"
)

func Example() {

	// Create a Hilbert curve for mapping to and from a 16 by 16 space.
	s, _ := hilbert.New(16)

	// Now map one dimension numbers in the range [0, N*N-1], to an x,y
	// coordinate on the curve where both x and y are in the range [0, N-1].
	x, y, _ := s.Map(96)

	// Also map back from (x,y) to t.
	t, _ := s.MapInverse(x, y)

	fmt.Printf("x = %d, y = %d, t = %d\n", x, y, t)

	// Output:
	// x = 4, y = 12, t = 96
}
