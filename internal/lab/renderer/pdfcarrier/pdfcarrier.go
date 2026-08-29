package pdfcarrier

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"unicode/utf8"

	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const RendererID = "pdf-carrier"
const VersionV1 = "v1"

type Renderer struct{}

func (Renderer) ID() string      { return RendererID }
func (Renderer) Version() string { return VersionV1 }

func (Renderer) Render(request renderer.Request) (renderer.Output, error) {
	if request.Width < 320 || request.Height < 320 {
		return renderer.Output{}, fmt.Errorf("pdf-carrier: dimensions must be >=320, got %dx%d", request.Width, request.Height)
	}
	if len(request.Source) == 0 || len(request.PageIndex) == 0 {
		return renderer.Output{}, fmt.Errorf("pdf-carrier: source and page index cannot be empty")
	}
	for _, page := range request.PageIndex {
		if page.ByteOffset < 0 || page.ByteLength < 0 || page.ByteOffset+page.ByteLength > len(request.Source) {
			return renderer.Output{}, fmt.Errorf("pdf-carrier: page %d has invalid byte range", page.Page)
		}
	}

	canvas := image.NewRGBA(image.Rect(0, 0, request.Width, request.Height))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.RGBA{R: 248, G: 247, B: 242, A: 255}), image.Point{}, draw.Src)
	margin := 24
	headerHeight := 40
	columns := ceilSqrt(len(request.PageIndex))
	rows := (len(request.PageIndex) + columns - 1) / columns
	tileWidth := (request.Width - 2*margin - (columns-1)*12) / columns
	tileHeight := (request.Height - 2*margin - headerHeight - (rows-1)*12) / rows
	if tileWidth < 100 || tileHeight < 80 {
		return renderer.Output{}, fmt.Errorf("pdf-carrier: %d pages do not fit requested canvas", len(request.PageIndex))
	}
	drawText(canvas, margin, margin+14, fmt.Sprintf("PDF CARRIER  |  %d pages  |  %d canonical bytes", len(request.PageIndex), len(request.Source)), canvas.Bounds())
	truncatedTiles := 0
	lostBytes := 0
	for pageIndex, page := range request.PageIndex {
		column := pageIndex % columns
		row := pageIndex / columns
		tile := image.Rect(margin+column*(tileWidth+12), margin+headerHeight+row*(tileHeight+12), margin+column*(tileWidth+12)+tileWidth, margin+headerHeight+row*(tileHeight+12)+tileHeight)
		draw.Draw(canvas, tile, image.NewUniform(color.White), image.Point{}, draw.Src)
		strokeRect(canvas, tile, color.RGBA{R: 55, G: 65, B: 81, A: 255})
		drawText(canvas, tile.Min.X+8, tile.Min.Y+15, fmt.Sprintf("P%d  @%d +%d", page.Page, page.ByteOffset, page.ByteLength), tile)
		coverage := drawWrapped(canvas, tile.Min.X+8, tile.Min.Y+31, page.Text, tile.Inset(8))
		if coverage.LostBytes > 0 {
			truncatedTiles++
			lostBytes += coverage.LostBytes
			drawText(canvas, tile.Min.X+8, tile.Max.Y-9, fmt.Sprintf("UNKNOWN: %d BYTES TRUNCATED", coverage.LostBytes), tile.Inset(8))
		}
	}

	pngBytes, err := renderer.EncodePNG(canvas)
	if err != nil {
		return renderer.Output{}, err
	}
	coverageStatus := "KNOWN"
	if truncatedTiles > 0 {
		coverageStatus = "UNKNOWN"
	}
	result := renderer.Result{Schema: renderer.ResultSchemaV1, RendererID: RendererID, Version: VersionV1, Width: request.Width, Height: request.Height, PixelSHA256: renderer.PixelHash(canvas), PNGSHA256: renderer.PNGHash(pngBytes), PNGBytes: len(pngBytes), SourceSHA256: request.SourceSHA, LayoutSeed: request.LayoutSeed, ImagePath: "images/original.png", CoverageStatus: coverageStatus, TruncatedTiles: truncatedTiles, LostBytes: lostBytes}
	return renderer.Output{PNG: pngBytes, Image: canvas, Result: result}, nil
}

type coverage struct{ LostBytes int }

func drawWrapped(canvas draw.Image, startX, startY int, text string, bounds image.Rectangle) coverage {
	const characterWidth = 7
	const lineHeight = 14
	maxCharacters := max(1, bounds.Dx()/characterWidth)
	lineY := startY
	consumedBytes := 0
	for _, paragraph := range strings.Split(text, "\n") {
		words := strings.Fields(paragraph)
		line := ""
		for wordIndex, word := range words {
			candidate := word
			if line != "" {
				candidate = line + " " + word
			}
			if utf8.RuneCountInString(candidate) > maxCharacters && line != "" {
				drawText(canvas, startX, lineY, line, bounds)
				lineY += lineHeight
				line = word
			} else {
				line = candidate
			}
			if lineY+lineHeight > bounds.Max.Y {
				return coverage{LostBytes: max(0, len(text)-consumedBytes)}
			}
			consumedBytes += len(word)
			if wordIndex > 0 {
				consumedBytes++
			}
		}
		if line != "" {
			drawText(canvas, startX, lineY, line, bounds)
			lineY += lineHeight
		}
		if lineY+lineHeight > bounds.Max.Y {
			return coverage{LostBytes: max(0, len(text)-consumedBytes)}
		}
		if consumedBytes < len(text) {
			consumedBytes++
		}
	}
	return coverage{}
}

func drawText(canvas draw.Image, startX, baselineY int, text string, bounds image.Rectangle) {
	drawer := &font.Drawer{Dst: canvas, Src: image.NewUniform(color.RGBA{R: 17, G: 24, B: 39, A: 255}), Face: basicfont.Face7x13, Dot: fixed.P(startX, baselineY)}
	for _, character := range text {
		if drawer.Dot.X.Ceil()+7 > bounds.Max.X {
			break
		}
		drawer.DrawString(string(character))
	}
}

func strokeRect(canvas draw.Image, rectangle image.Rectangle, stroke color.Color) {
	draw.Draw(canvas, image.Rect(rectangle.Min.X, rectangle.Min.Y, rectangle.Max.X, rectangle.Min.Y+1), image.NewUniform(stroke), image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(rectangle.Min.X, rectangle.Max.Y-1, rectangle.Max.X, rectangle.Max.Y), image.NewUniform(stroke), image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(rectangle.Min.X, rectangle.Min.Y, rectangle.Min.X+1, rectangle.Max.Y), image.NewUniform(stroke), image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(rectangle.Max.X-1, rectangle.Min.Y, rectangle.Max.X, rectangle.Max.Y), image.NewUniform(stroke), image.Point{}, draw.Src)
}

func ceilSqrt(number int) int {
	root := 1
	for root*root < number {
		root++
	}
	return root
}
