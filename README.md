# ASCIIPIC

ASCIIPIC is a high-performance command-line tool written in Go that converts images into ASCII art. It supports true-color terminal output and includes a specialized PNG exporter that applies a background luminosity effect, producing vibrant, high-contrast results suitable for sharing.

---

## Features

### True Color Support

Renders ASCII art in the terminal using 24-bit ANSI color codes, preserving the original image's color palette.

### High-Quality PNG Export

Exports ASCII art to PNG with a custom glow rendering engine. This engine fills the background of each character cell with a dimmed version of the image color, preventing exported art from appearing too dark or washed out.

### Multiple Character Sets

Includes 15 character sets ranging from classic ASCII to Braille patterns, block elements, and mathematical symbols.

### Detailed Configuration

Complete control over output width, font size for PNG export, and brightness adjustment.

### Broad Format Support

Supports JPEG and PNG input files.

---

## Installation

### Prerequisites

- Go 1.18 or higher

### Install via Go

```sh
go install github.com/thesfb/asciipic@latest
```

### Build from Source

```sh
git clone https://github.com/thesfb/asciipic.git
cd asciipic
go build -o asciipic main.go
```

---

## Usage

### Basic Terminal Output

```sh
./asciipic -input image.jpg
```

### True Color Terminal Output

```sh
./asciipic -input image.jpg -color -width 100
```

### Export to PNG

```sh
./asciipic -input image.jpg -png -color -output result.png -width 200
```

### Advanced Example

Using the Braille character set with enhanced brightness:

```sh
./asciipic -input image.jpg -png -color -charset e -width 300 -brightness 130 -output braille_art.png
```

---

## Command Line Flags

| Flag        | Type   | Default          | Description                                                    |
|-------------|--------|------------------|----------------------------------------------------------------|
| `-input`    | string | `""`             | Path to the input image file (required).                       |
| `-output`   | string | `ascii-art.png`  | Path for the exported PNG file.                               |
| `-width`    | uint   | `80`             | Width of the output in characters.                            |
| `-color`    | bool   | `false`          | Enable true color output (ANSI for terminal, RGB for PNG).    |
| `-png`      | bool   | `false`          | Enable PNG export mode.                                       |
| `-charset`  | string | `a`              | Select character set (a–o).                                   |
| `-fontsize` | int    | `8`              | Font size in pixels for PNG export.                           |
| `-brightness` | int  | `110`            | Brightness percentage (50–200).                               |

---

## Character Sets

- `a`: Classic ASCII  
- `b`: Blocks  
- `c`: Dots  
- `d`: Vertical Bars  
- `e`: Braille  
- `f`: Numbers  
- `g`: Letters  
- `h`: Slashes  
- `i`: Punctuation  
- `j`: Mixed Blocks  
- `k`: Advanced Braille  
- `l`: Circle Fill  
- `m`: Box Drawing  
- `n`: Dot Sizes  
- `o`: Math Symbols  

---

## Dependencies

- `github.com/golang/freetype`
- `github.com/nfnt/resize`
- `golang.org/x/image`

---

## License

Distributed under the MIT License. See `LICENSE` for details.
