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

package lib

import (
	"image"
	"image/color"
	"image/draw"
)

// UniqueColors returns the first 256 unique color.Color used in this image.
func UniqueColors(src image.Image) []color.Color {
	var colors []color.Color

	bounds := src.Bounds()

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			c := src.At(x, y)
			found := false
			for i := 0; i < len(colors) && !found; i++ {
				if colors[i] == c {
					found = true
				}
			}
			if !found {
				colors = append(colors, c)
				if len(colors) >= 256 {
					return colors
				}
			}
		}
	}

	return colors
}

// ConvertToPaletted converts the given image into a paletted one.
// Colors are converted using a naive approache. The first 256 unique colors
// are retained, and the rest are mapped to hopefully a nearby color.
func ConvertToPaletted(src image.Image) *image.Paletted {

	if dst, ok := src.(*image.Paletted); ok {
		return dst
	}

	bounds := src.Bounds()
	colors := UniqueColors(src)

	dst := image.NewPaletted(bounds, color.Palette(colors))
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}
