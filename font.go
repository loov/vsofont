// package vsofont implements parsing font format in http://www.villehelin.com/vsofont.html
//
//    # some text, perhaps explanations
//
//    JUMP!
//
//    GRID: <GRID X> x <GRID Y>
//    SPACING: <EMPTY SPACE BETWEEN THE CHARACTERS>
//    SCALING: <SCALING X> x <SCALING Y>
//    COLOR: <R> <G> <B> <A>
//
//    <CHARACTER> <LINES, DEFINED USING INDICES TO THE GRID> -1
//    ...

package vsofont

import (
	"fmt"
	"strconv"
	"strings"
)

type Font struct {
	Spacing float32
	Glyphs  map[string]Glyph
}

type Glyph struct {
	Rune  string
	Lines [][2]Vector
}

type Vector = struct {
	X, Y float32
}

func MustParse(text string) *Font {
	font, err := Parse(text)
	if err != nil {
		panic(err)
	}
	return font
}

func Parse(text string) (*Font, error) {
	font := &Font{
		Glyphs:  map[string]Glyph{},
		Spacing: 0,
	}

	var gridWidth, gridHeight int
	_ = gridHeight
	var scaleWidth, scaleHeight float32

	convert := func(coord int) Vector {
		return Vector{
			X: float32(coord%gridWidth) * scaleWidth,
			Y: float32(coord/gridWidth) * scaleHeight,
		}
	}

	jumpFound := false
	for lineno, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		if line == "JUMP!" {
			jumpFound = true
			continue
		} else if !jumpFound {
			continue
		}

		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "GRID:":
			// GRID: 5 x 5
			if len(tokens) != 4 {
				return nil, fmt.Errorf("invalid number of tokens")
			}
			width, err := strconv.Atoi(tokens[1])
			if err != nil {
				return nil, fmt.Errorf("failed to grid width: %v", err)
			}

			height, err := strconv.Atoi(tokens[3])
			if err != nil {
				return nil, fmt.Errorf("failed to grid height: %v", err)
			}

			gridWidth, gridHeight = width, height

		case "SPACING:":
			if len(tokens) != 2 {
				return nil, fmt.Errorf("invalid number of tokens")
			}
			// SPACING: 0.005
			spacing, err := strconv.ParseFloat(tokens[1], 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse spacing: %v", err)
			}
			font.Spacing = float32(spacing)

		case "SCALING:":
			if len(tokens) != 4 {
				return nil, fmt.Errorf("invalid number of tokens")
			}
			// SCALING: 0.2 x 0.2
			width, err := strconv.ParseFloat(tokens[1], 32)
			if err != nil {
				return nil, fmt.Errorf("failed to scaling width: %v", err)
			}

			height, err := strconv.ParseFloat(tokens[3], 32)
			if err != nil {
				return nil, fmt.Errorf("failed to scaling height: %v", err)
			}

			scaleWidth, scaleHeight = float32(width), float32(height)
		case "COLOR:":
			// ignore
			continue
		default:
			glyph := Glyph{
				Rune: tokens[0],
			}

			for i := 1; i < len(tokens); i += 2 {
				if i+1 >= len(tokens) {
					break
				}

				a, err := strconv.Atoi(tokens[i])
				if err != nil {
					return nil, fmt.Errorf("%d: failed to read %q: %v", lineno, tokens[i], err)
				}

				b, err := strconv.Atoi(tokens[i+1])
				if err != nil {
					return nil, fmt.Errorf("%d: failed to read %q: %v", lineno, tokens[i+1], err)
				}

				glyph.Lines = append(glyph.Lines, [2]Vector{convert(a), convert(b)})
			}

			font.Glyphs[glyph.Rune] = glyph
		}

	}

	return font, nil
}
