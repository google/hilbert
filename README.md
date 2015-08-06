# Hilbert

Golang package for mapping values to and from a Hilbert curve.

![Image of 8 by 8 Hilbert curve](images/hilbert.png)

Converted from the code available on Wikipedia, with additional help from:
 * https://en.wikipedia.org/wiki/Hilbert_curve
 * http://bit-player.org/2013/mapping-the-hilbert-curve

[Documentation available here](https://godoc.org/github.com/google/hilbert)

*This is not an official Google product (experimental or otherwise), it is just code that happens to be owned by Google.*
 
## How to use

Install:

```bash
go get github.com/google/hilbert
```

Example:

```go
import "github.com/google/hilbert"
	
// Create a Hilbert curve for mapping to and from a 16 by 16 space.
s, err := hilbert.New(16)

// Now map one dimension numbers in the range [0, N*N-1], to an x,y
// coordinate on the curve where both x and y are in the range [0, N-1].
x, y, err := s.Map(t)

// Also map back from (x,y) to t.
t, err := s.MapInverse(x, y)
```

## Demo

The demo directory contains an example on how to draw an image of a Hilbert
curve, as well as animations of varying sized Hilbert spaces.

```bash
go run $GOPATH/src/github.com/google/hilbert/demo/demo.go
```

and the following two images are generated. One a simple 8x8 grid:

![Image of 8 by 8 Hilbert curve](images/hilbert.png)

and an animation of N in the range 1..8:

![Hilbert curve animation](images/hilbert_animation.gif)

## Licence (Apache 2)

```
Copyright 2015 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
