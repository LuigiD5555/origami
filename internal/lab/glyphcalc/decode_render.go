package glyphcalc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// DecodeGlyphGridPNG decodes a grid produced by RenderGlyphGrid without a
// semantic sidecar. The layout parameters are transport metadata, not carrier
// semantic meaning. Trailing all-white cells are ignored.
//
// cell >= 56 is required so every current glyph feature remains inside its own
// cell; this keeps cells independently hash-decodable.
func DecodeGlyphGridPNG(pngBytes []byte, columns, cell, margin int) ([]Glyph, error) {
	if columns <= 0 { columns = 24 }
	if cell < 56 { return nil, fmt.Errorf("decoder requires cell >= 56, got %d", cell) }
	if margin < 4 { margin = 4 }
	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil { return nil, fmt.Errorf("decode png: %w", err) }
	bounds := img.Bounds()
	expectedWidth := margin*2 + columns*cell
	if bounds.Dx() != expectedWidth {
		return nil, fmt.Errorf("unexpected grid width: got=%d want=%d", bounds.Dx(), expectedWidth)
	}
	usableHeight := bounds.Dy() - 2*margin
	if usableHeight < 0 || usableHeight%cell != 0 {
		return nil, fmt.Errorf("unexpected grid height %d for cell=%d margin=%d", bounds.Dy(), cell, margin)
	}
	rows := usableHeight / cell

	templates, blank, err := templateHashes(cell)
	if err != nil { return nil, err }
	out := make([]Glyph, 0, rows*columns)
	trailingBlank := false
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			x0 := bounds.Min.X + margin + col*cell
			y0 := bounds.Min.Y + margin + row*cell
			h := cellHash(img, image.Rect(x0, y0, x0+cell, y0+cell))
			if h == blank {
				trailingBlank = true
				continue
			}
			if trailingBlank {
				return nil, fmt.Errorf("non-blank glyph after trailing blank cell at row=%d col=%d", row, col)
			}
			value, ok := templates[h]
			if !ok {
				return nil, fmt.Errorf("unrecognized or corrupted glyph at row=%d col=%d", row, col)
			}
			out = append(out, EncodeByte(len(out), value))
		}
	}
	if len(out) == 0 { return nil, fmt.Errorf("no glyphs decoded") }
	return out, nil
}

func templateHashes(cell int) (map[[32]byte]byte, [32]byte, error) {
	whiteCell := image.NewRGBA(image.Rect(0, 0, cell, cell))
	draw.Draw(whiteCell, whiteCell.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	blank := cellHash(whiteCell, whiteCell.Bounds())

	out := make(map[[32]byte]byte, 256)
	for i := 0; i < 256; i++ {
		img := image.NewRGBA(image.Rect(0, 0, cell, cell))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
		drawGlyph(img, cell/2, cell/2, cell, EncodeByte(0, byte(i)))
		h := cellHash(img, img.Bounds())
		if prior, exists := out[h]; exists {
			return nil, blank, fmt.Errorf("glyph rendering collision: byte %d and %d", prior, i)
		}
		out[h] = byte(i)
	}
	return out, blank, nil
}

func cellHash(img image.Image, rect image.Rectangle) [32]byte {
	buf := make([]byte, 0, rect.Dx()*rect.Dy()*4)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
	return sha256.Sum256(buf)
}
