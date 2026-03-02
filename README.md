# asciizer

A command-line tool that converts images to ASCII art. Supports JPEG, PNG, and GIF with automatic format detection.

## Features

- Multi-format support (JPEG, PNG, GIF)
- Automatic resizing with aspect ratio correction for terminal display
- ANSI 256-color output
- Configurable character ramps (10-char short or 70-char detailed)
- Invertible brightness mapping
- File or stdout output

## Install

```
make deploy
```

This builds the binary, installs it to `~/.local/bin/`, installs the man page, and sets up zsh completions.

## Usage

```
asciizer [flags] <image_file>
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-w` | 80 | Output width in characters (0 = no resize) |
| `-o` | `<input>.asc` | Output file path |
| `-stdout` | false | Print to stdout instead of file |
| `-invert` | false | Reverse brightness mapping |
| `-color` | false | ANSI 256-color output |
| `-full-ramp` | false | Use 70-char gradient instead of 10-char |

### Examples

```bash
# Convert and save to photo.asc
asciizer photo.jpg

# Print to terminal
asciizer -stdout photo.png

# Color output at 120 columns
asciizer -stdout -color -w 120 photo.png

# Inverted with detailed ramp
asciizer -stdout -invert -full-ramp photo.gif

# Original resolution, custom output path
asciizer -o art.txt -w 0 photo.jpg
```

## Build

```
make build
```

## License

MIT
