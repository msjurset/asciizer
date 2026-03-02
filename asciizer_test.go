package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// newSolidImage creates a uniform image of the given color and dimensions.
func newSolidImage(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestResizeImage(t *testing.T) {
	src := newSolidImage(200, 100, color.White)

	resized := resizeImage(src, 40)
	bounds := resized.Bounds()

	if bounds.Dx() != 40 {
		t.Errorf("expected width 40, got %d", bounds.Dx())
	}
	if bounds.Dy() < 1 {
		t.Errorf("expected height >= 1, got %d", bounds.Dy())
	}

	// Pixel values should survive: white stays white.
	r, g, b, a := resized.At(0, 0).RGBA()
	if r != 0xffff || g != 0xffff || b != 0xffff || a != 0xffff {
		t.Errorf("expected white pixel, got (%d,%d,%d,%d)", r, g, b, b)
	}
}

func TestResizeImageMinHeight(t *testing.T) {
	// Very wide, very short image should clamp to height 1.
	src := newSolidImage(1000, 1, color.Black)
	resized := resizeImage(src, 50)
	if resized.Bounds().Dy() < 1 {
		t.Errorf("expected height >= 1, got %d", resized.Bounds().Dy())
	}
}

func TestToGrayscale(t *testing.T) {
	img := newSolidImage(2, 2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	gray := toGrayscale(img)

	if gray.Bounds().Dx() != 2 || gray.Bounds().Dy() != 2 {
		t.Fatalf("expected 2x2 gray, got %dx%d", gray.Bounds().Dx(), gray.Bounds().Dy())
	}

	// Red luminance: 0.299*255 ≈ 76
	val := gray.GrayAt(0, 0).Y
	if val < 70 || val > 82 {
		t.Errorf("expected luminance ~76 for pure red, got %d", val)
	}
}

func TestToGrayscaleWhite(t *testing.T) {
	img := newSolidImage(1, 1, color.White)
	gray := toGrayscale(img)
	val := gray.GrayAt(0, 0).Y
	if val < 250 {
		t.Errorf("expected luminance ~255 for white, got %d", val)
	}
}

func TestToGrayscaleBlack(t *testing.T) {
	img := newSolidImage(1, 1, color.Black)
	gray := toGrayscale(img)
	val := gray.GrayAt(0, 0).Y
	if val > 5 {
		t.Errorf("expected luminance ~0 for black, got %d", val)
	}
}

func TestGenerateASCII(t *testing.T) {
	img := newSolidImage(3, 2, color.White)
	gray := toGrayscale(img)

	art := generateASCII(img, gray, rampShort, false, false)
	lines := strings.Split(strings.TrimRight(art, "\n"), "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if len(lines[0]) != 3 {
		t.Errorf("expected 3 chars per line, got %d", len(lines[0]))
	}

	// White maps to the last char in the ramp (lightest).
	lastChar := rampShort[len(rampShort)-1]
	for _, ch := range lines[0] {
		if byte(ch) != lastChar {
			t.Errorf("expected '%c' for white pixel, got '%c'", lastChar, ch)
		}
	}
}

func TestGenerateASCIIInvert(t *testing.T) {
	img := newSolidImage(1, 1, color.White)
	gray := toGrayscale(img)

	normal := generateASCII(img, gray, rampShort, false, false)
	inverted := generateASCII(img, gray, rampShort, true, false)

	normalChar := strings.TrimSpace(normal)
	invertedChar := strings.TrimSpace(inverted)

	if normalChar == invertedChar {
		t.Errorf("invert should produce a different character, both got '%s'", normalChar)
	}
}

func TestGenerateASCIIColor(t *testing.T) {
	img := newSolidImage(2, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	gray := toGrayscale(img)

	art := generateASCII(img, gray, rampShort, false, true)

	if !strings.Contains(art, "\033[38;5;") {
		t.Error("color mode should contain ANSI escape codes")
	}
	if !strings.Contains(art, "\033[0m") {
		t.Error("color mode should contain reset sequence")
	}
}

func TestGenerateASCIIFullRamp(t *testing.T) {
	img := newSolidImage(1, 1, color.Black)
	gray := toGrayscale(img)

	artShort := generateASCII(img, gray, rampShort, false, false)
	artFull := generateASCII(img, gray, rampStandard, false, false)

	// Black maps to first char in each ramp, which differ.
	shortChar := strings.TrimSpace(artShort)
	fullChar := strings.TrimSpace(artFull)

	if shortChar != string(rampShort[0]) {
		t.Errorf("short ramp: expected '%c' for black, got '%s'", rampShort[0], shortChar)
	}
	if fullChar != string(rampStandard[0]) {
		t.Errorf("full ramp: expected '%c' for black, got '%s'", rampStandard[0], fullChar)
	}
}
