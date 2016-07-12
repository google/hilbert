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

// Package main is a simple demo to show how to use the hilbert library
// When ran, this demo will create two images, hilbert.png, and
// hilbert_animation.gif.
//
// It is suggested you optimise/compress both images before uploading.
//     optipng hilbert.png
//     gifsicle -O -o hilbert_animation_compressed.gif hilbert_animation.gif
//
package main

import (
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/hilbert"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

// HilbertImage facilitates the drawing of a Hilbert Curve
type HilbertImage struct {
	N          int     // Hilbert Space is N by N
	SquareSize float64 // Size of each square in pixels

	DrawText   bool    // Should text be drawn on the image
	TextMargin float64 // Margin around text in pixels

	BackgroundColor color.Color
	GridColor       color.Color
	TextColor       color.Color
	SnakeColor      color.Color
}

func (h *HilbertImage) toPixel(x, y int) (px, py float64) {
	return float64(x) * h.SquareSize, float64(y) * h.SquareSize
}

func (h *HilbertImage) createImage() (draw.Image, error) {
	width, height := h.toPixel(h.N, h.N)
	return image.NewRGBA(image.Rect(0, 0, int(width), int(height))), nil
}

func (h *HilbertImage) drawRectange(gc draw2d.GraphicContext, px1, py1, px2, py2 float64) {
	gc.SetFillColor(h.BackgroundColor)
	gc.SetStrokeColor(h.GridColor)
	gc.SetLineWidth(1)

	draw2dkit.Rectangle(gc, px1, py1, px2, py2)
	gc.FillStroke()
}

func (h *HilbertImage) drawText(gc draw2d.GraphicContext, px1, py1 float64, t int) {
	if !h.DrawText {
		return
	}

	text := strconv.Itoa(t)
	_, top, _, _ := gc.GetStringBounds(text)

	gc.SetFillColor(h.TextColor)
	gc.FillStringAt(text, px1+h.TextMargin, py1-top+h.TextMargin)
}

func (h *HilbertImage) drawSnake(gc draw2d.GraphicContext, snake *draw2d.Path) {
	gc.SetStrokeColor(h.SnakeColor)
	gc.SetLineCap(draw2d.SquareCap)
	gc.SetLineJoin(draw2d.MiterJoin)
	gc.SetLineWidth(2)

	gc.Stroke(snake)
}

// CreateHilbertImage returns a new hilbertImage ready for drawing.
func CreateHilbertImage(n int, squareSize float64) *HilbertImage {
	return &HilbertImage{
		N:          n,
		SquareSize: squareSize,

		DrawText:   true,
		TextMargin: 3.0,

		BackgroundColor: color.RGBA{0xee, 0xee, 0xff, 0xff},
		GridColor:       color.White,
		TextColor:       color.RGBA{0x33, 0x33, 0x33, 0xff},
		SnakeColor:      color.RGBA{0x33, 0x33, 0x33, 0xff},
	}
}

// Draw uses the parameters in the hilbertImage and returns a Image
func (h *HilbertImage) Draw() (draw.Image, error) {

	// Create a Hilbert Curve Mapper
	s, err := hilbert.New(h.N)
	if err != nil {
		return nil, fmt.Errorf("Failed to create hilbert space: %s", err.Error())
	}

	img, err := h.createImage()
	if err != nil {
		return nil, err
	}

	gc := draw2dimg.NewGraphicContext(img)
	snake := &draw2d.Path{}

	for t := 0; t < h.N*h.N; t++ {

		// Map the 1D number into the 2D space
		x, y, err := s.Map(t)
		if err != nil {
			return nil, err
		}

		px1, py1 := h.toPixel(x, y)
		px2, py2 := h.toPixel(x+1, y+1)

		// Draw the grid for t
		h.drawRectange(gc, px1, py1, px2, py2)
		h.drawText(gc, px1, py1, t)

		// Move the snake along
		centerX, centerY := px1+h.SquareSize/2, py1+h.SquareSize/2
		if t == 0 {
			snake.MoveTo(centerX, centerY)
		} else {
			snake.LineTo(centerX, centerY)
		}
	}

	// Draw the snake at the end, to form one continuous line.
	h.drawSnake(gc, snake)

	return img, nil
}

// uniqueColors returns the first 256 unique color.Color used in this image.
// TODO consider moving into a helper/graphics library
func uniqueColors(src image.Image) []color.Color {
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

// convertToPaletted converts the given image into a paletted one.
// Colors are converted using a naive approache. The first 256 unique colors
// are retained, and the rest are mapped to hopefully a nearby color.
func convertToPaletted(src image.Image) *image.Paletted {

	if dst, ok := src.(*image.Paletted); ok {
		return dst
	}

	bounds := src.Bounds()
	colors := uniqueColors(src)

	dst := image.NewPaletted(bounds, color.Palette(colors))
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// setupDraw2D ensure Draw2D can find its fonts.
func setupDraw2D() {
	p, err := build.Default.Import("github.com/llgcode/draw2d", "", build.FindOnly)
	if err != nil {
		log.Fatalf("Couldn't find llgcode/draw2d files: %v", err)
	}

	draw2d.SetFontFolder(filepath.Join(p.Dir, "resource", "font"))
}

func mainDrawOne() error {
	log.Printf("Drawing one image")

	img, err := CreateHilbertImage(8, 64).Draw()
	if err != nil {
		return err
	}
	return draw2dimg.SaveToPngFile("hilbert.png", img)
}

func mainDrawAnimation() error {
	log.Printf("Drawing animation")

	iterations := 8
	imageWidth := 512.0

	g := gif.GIF{
		Image:     make([]*image.Paletted, iterations),
		Delay:     make([]int, iterations),
		LoopCount: 0,
	}

	for i := 0; i < iterations; i++ {
		log.Printf("    Drawing frame %d", i)

		n := 1 << uint(i+1)
		h := CreateHilbertImage(n, imageWidth/float64(n))
		h.DrawText = false
		img, err := h.Draw()
		if err != nil {
			return err
		}

		g.Image[i] = convertToPaletted(img) // draw2d doesn't support paletted images, so we convert
		g.Delay[i] = 200                    // 200 x 100th of a second = 2 second
	}

	f, err := os.Create("hilbert_animation.gif")
	if err != nil {
		return err
	}
	return gif.EncodeAll(f, &g)
}

func main() {

	setupDraw2D()

	err := mainDrawOne()
	if err != nil {
		log.Fatalf("Failed to draw image: %s", err.Error())
	}

	err = mainDrawAnimation()
	if err != nil {
		log.Fatalf("Failed to draw animation: %s", err.Error())
	}

}
