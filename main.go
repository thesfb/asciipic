package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/gofont/gomonobold"
)

// Character sets for different styles
var charSets = map[string][]rune{
	"a": {' ', '.', ':', '-', '=', '+', '*', '#', '%', '@'},
	"b": {' ', '░', '▒', '▓', '█'},
	"c": {' ', '·', '•', '○', '◉', '●'},
	"d": {' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'},
	"e": {' ', '⣀', '⣄', '⣤', '⣦', '⣶', '⣷', '⣿'},
	"f": {' ', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
	"g": {' ', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'},
	"h": {' ', '/', '\\', '|', '-', '+', 'x', '*', '#', '@'},
	"i": {' ', '`', '.', ',', ':', ';', '!', '>', '<', '~'},
	"j": {' ', '░', '▒', '▓', '█', '▀', '▄', '▌', '▐'},
	"k": {'⠀', '⠁', '⠃', '⠇', '⠏', '⠟', '⠿', '⡿', '⣿'},
	"l": {' ', '○', '◔', '◐', '◕', '⬤'},
	"m": {' ', '┤', '┴', '├', '┬', '┼', '╬', '█'},
	"n": {' ', '˙', '·', '•', '●', '⚫'},
	"o": {' ', '⋅', '∘', '∙', '○', '◎', '⦿', '●'},
}

type Config struct {
	Width      uint
	Color      bool
	Input      string
	CharSet    string
	FontSize   int
	ExportPNG  bool
	OutputPath string
	Brightness int
}

func getCharSet(setName string) []rune {
	if chars, ok := charSets[setName]; ok {
		return chars
	}
	return charSets["a"] // this is thedefault charset
}

func brightnessToASCII(brightness uint8, chars []rune) rune {
	idx := int(brightness) * len(chars) / 256
	if idx >= len(chars) {
		idx = len(chars) - 1
	}
	return chars[idx]
}

func rgbToBrightness(r, g, b uint32) uint8 {
	return uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func convertImage(img image.Image, width uint, useColor bool, charSet string) string {
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())
	height := uint(float64(width) * aspectRatio * 0.5)

	resized := resize.Resize(width, height, img, resize.Lanczos3)
	var output bytes.Buffer

	chars := getCharSet(charSet)

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			brightness := rgbToBrightness(r, g, b)
			ch := brightnessToASCII(brightness, chars)

			if useColor {
				r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

				output.WriteString(fmt.Sprintf(
					"\033[38;2;%d;%d;%dm%c\033[0m",
					r8, g8, b8, ch,
				))
			} else {
				output.WriteRune(ch)
			}
		}
		output.WriteRune('\n')
	}

	return output.String()
}

func convertImageToPNG(img image.Image, width uint, useColor bool, charSet string, fontSize int, outputPath string, brightnessBoost int) error {
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())
	height := uint(float64(width) * aspectRatio * 0.5)

	resized := resize.Resize(width, height, img, resize.Lanczos3)
	chars := getCharSet(charSet)

	charWidth := fontSize
	charHeight := int(float64(fontSize) * 1.8)
	padding := 40
	outWidth := int(width)*charWidth + padding*2
	outHeight := int(height)*charHeight + padding*2

	outImg := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))

	bgColor := color.RGBA{R: 5, G: 5, B: 5, A: 255}
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	font, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(float64(fontSize))
	c.SetClip(outImg.Bounds())
	c.SetDst(outImg)

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			brightness := rgbToBrightness(r, g, b)
			ch := brightnessToASCII(brightness, chars)

			rVal := int(r >> 8)
			gVal := int(g >> 8)
			bVal := int(b >> 8)

			r8 := uint8(min(rVal*brightnessBoost/100, 255))
			g8 := uint8(min(gVal*brightnessBoost/100, 255))
			b8 := uint8(min(bVal*brightnessBoost/100, 255))

			bgR := uint8(min(rVal*30/100, 255))
			bgG := uint8(min(gVal*30/100, 255))
			bgB := uint8(min(bVal*30/100, 255))

			posX := x*charWidth + padding
			posY := y*charHeight + padding

			if useColor {

				cellRect := image.Rect(posX, posY, posX+charWidth, posY+charHeight)
				draw.Draw(outImg, cellRect, &image.Uniform{color.RGBA{bgR, bgG, bgB, 255}}, image.Point{}, draw.Src)

				c.SetSrc(&image.Uniform{color.RGBA{r8, g8, b8, 255}})
			} else {

				grayLevel := uint8(min(int(brightness)*brightnessBoost/100, 255))
				c.SetSrc(&image.Uniform{color.RGBA{grayLevel, grayLevel, grayLevel, 255}})
			}

			pt := freetype.Pt(posX, posY+int(c.PointToFixed(float64(fontSize))>>6))
			c.DrawString(string(ch), pt)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create the file: %w", err)
	}
	defer outFile.Close()

	return png.Encode(outFile, outImg)
}

func convertFile(config Config) error {
	file, err := os.Open(config.Input)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if config.ExportPNG {
		outputPath := config.OutputPath
		if outputPath == "" {
			outputPath = "ascii-art.png"
		}

		brightnessBoost := config.Brightness
		if brightnessBoost == 0 {
			brightnessBoost = 110
		}

		err := convertImageToPNG(img, config.Width, config.Color, config.CharSet, config.FontSize, outputPath, brightnessBoost)
		if err != nil {
			return fmt.Errorf("failed to export PNG: %w", err)
		}
		fmt.Printf("Exported to: %s (brightness: %d%%)\n", outputPath, brightnessBoost)
		return nil
	}

	// Standard Terminal Output
	ascii := convertImage(img, config.Width, config.Color, config.CharSet)
	fmt.Print(ascii)
	return nil
}

func printHelp() {
	fmt.Println("ASCIIPIC - Image to ASCII Art Converter")
	fmt.Println("\nAvailable Character Sets:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	descriptions := map[string]string{
		"a": "Classic ASCII   ",
		"b": "Blocks          ",
		"c": "Dots            ",
		"d": "Vertical Bars   ",
		"e": "Braille         ",
		"f": "Numbers         ",
		"g": "Letters         ",
		"h": "Slashes         ",
		"i": "Punctuation     ",
		"j": "Mixed Blocks    ",
		"k": "Braille Advanced",
		"l": "Circle Fill     ",
		"m": "Box Drawing     ",
		"n": "Dot Sizes       ",
		"o": "Math Symbols    ",
	}

	for k := 'a'; k <= 'o'; k++ {
		key := string(k)
		chars := charSets[key]
		desc := descriptions[key]
		fmt.Printf("  %s: %s → %s\n", key, desc, string(chars))
	}

	fmt.Println("\nUsage Examples:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  # Basic terminal output")
	fmt.Println("  asciipic -input photo.jpg")
	fmt.Println()
	fmt.Println("  # With TRUE color (for terminal)")
	fmt.Println("  asciipic -input photo.jpg -color")
	fmt.Println()
	fmt.Println("  # Braille characters (high detail)")
	fmt.Println("  asciipic -input photo.jpg -charset e -color -width 150")
	fmt.Println()
	fmt.Println("  # Export as PNG image (Vibrant Glow effect)")
	fmt.Println("  asciipic -input photo.jpg -png -color -charset b -fontsize 12 -width 200 -output art.png")
	fmt.Println()
	fmt.Println("  # Darker/Moody PNG export")
	fmt.Println("  asciipic -input photo.jpg -png -color -brightness 80 -output dark.png")
	fmt.Println()
	fmt.Println(" Flags:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	flag.PrintDefaults()
}

func main() {
	var config Config

	flag.StringVar(&config.Input, "input", "", "Path to input image file")
	flag.UintVar(&config.Width, "width", 80, "Output width in characters")
	flag.BoolVar(&config.Color, "color", false, "Enable TRUE color output (preserves original colors)")
	flag.BoolVar(&config.ExportPNG, "png", false, "Export as png image")
	flag.StringVar(&config.OutputPath, "output", "ascii-art.png", "Output path for png export")
	flag.StringVar(&config.CharSet, "charset", "a", "Character set to use (a-o). Use --help to see all sets")
	flag.IntVar(&config.FontSize, "fontsize", 8, "Font size in pixels for PNG export (1-50)")
	flag.IntVar(&config.Brightness, "brightness", 110, "Brightness boost for PNG export (50-200, default 110 = +10%)")

	flag.Usage = printHelp

	flag.Parse()

	config.CharSet = strings.ToLower(config.CharSet)
	if _, ok := charSets[config.CharSet]; !ok {
		fmt.Printf("Invalid charset '%s'. Use --help to see available sets.\n", config.CharSet)
		os.Exit(1)
	}

	if config.Input != "" {
		if err := convertFile(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		printHelp()
		os.Exit(1)
	}
}
