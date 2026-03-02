package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	rampStandard = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	rampShort    = "@%#*+=-:. "
)

var version = "dev"

type options struct {
	input    string
	output   string
	width    int
	stdout   bool
	invert   bool
	color    bool
	fullRamp bool
}

func main() {
	width := flag.Int("w", 80, "output width in characters (0 = no resize)")
	output := flag.String("o", "", "output file path (default: <input>.asc)")
	stdout := flag.Bool("stdout", false, "print to stdout instead of file")
	invert := flag.Bool("invert", false, "reverse brightness mapping")
	colorFlag := flag.Bool("color", false, "ANSI 256-color output")
	fullRamp := flag.Bool("full-ramp", false, "use 70-char gradient instead of 10-char")
	showVersion := flag.Bool("version", false, "print version and exit")
	completion := flag.String("completion", "", "print shell completion script (zsh, bash)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "asciizer %s\n\nUsage: asciizer [flags] <image_file>\n\nFlags:\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("asciizer %s\n", version)
		return
	}

	if *completion != "" {
		switch *completion {
		case "zsh":
			fmt.Print(zshCompletion)
		case "bash":
			fmt.Print(bashCompletion)
		default:
			fmt.Fprintf(os.Stderr, "unsupported shell: %s (supported: zsh, bash)\n", *completion)
			os.Exit(1)
		}
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputFile := flag.Arg(0)
	outFile := *output
	if outFile == "" {
		ext := filepath.Ext(inputFile)
		outFile = strings.TrimSuffix(inputFile, ext) + ".asc"
	}

	opts := options{
		input:    inputFile,
		output:   outFile,
		width:    *width,
		stdout:   *stdout,
		invert:   *invert,
		color:    *colorFlag,
		fullRamp: *fullRamp,
	}

	if err := run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	f, err := os.Open(opts.input)
	if err != nil {
		return fmt.Errorf("opening image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decoding image: %w", err)
	}

	if opts.width > 0 {
		img = resizeImage(img, opts.width)
	}

	gray := toGrayscale(img)

	ramp := rampShort
	if opts.fullRamp {
		ramp = rampStandard
	}

	art := generateASCII(img, gray, ramp, opts.invert, opts.color)

	if opts.stdout {
		fmt.Print(art)
		return nil
	}

	if err := os.WriteFile(opts.output, []byte(art), 0644); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	fmt.Printf("Saved ASCII art to %s\n", opts.output)
	return nil
}

func resizeImage(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	aspectRatio := float64(srcH) / float64(srcW)
	targetHeight := int(math.Round(float64(targetWidth) * aspectRatio * 0.5))
	if targetHeight < 1 {
		targetHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	for y := 0; y < targetHeight; y++ {
		srcY := bounds.Min.Y + y*srcH/targetHeight
		for x := 0; x < targetWidth; x++ {
			srcX := bounds.Min.X + x*srcW/targetWidth
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			lum := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
			gray.SetGray(x, y, color.Gray{Y: lum})
		}
	}
	return gray
}

func generateASCII(img image.Image, gray *image.Gray, ramp string, invert bool, colorMode bool) string {
	bounds := gray.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rampLen := len(ramp)

	var b strings.Builder
	b.Grow(w*h*2 + h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			grayVal := gray.GrayAt(x, y).Y

			idx := int(float64(grayVal) * float64(rampLen-1) / 255.0)
			if invert {
				idx = rampLen - 1 - idx
			}

			ch := ramp[idx]

			if colorMode {
				r, g, bl, _ := img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y).RGBA()
				r8, g8, b8 := r>>8, g>>8, bl>>8
				code := 16 + (r8/51)*36 + (g8/51)*6 + b8/51
				fmt.Fprintf(&b, "\033[38;5;%dm%c", code, ch)
			} else {
				b.WriteByte(ch)
			}
		}
		if colorMode {
			b.WriteString("\033[0m")
		}
		b.WriteByte('\n')
	}
	return b.String()
}
