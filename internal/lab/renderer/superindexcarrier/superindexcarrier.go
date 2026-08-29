package superindexcarrier

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/LuigiD5555/origami/internal/lab/capacity"
	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"github.com/LuigiD5555/origami/internal/lab/superindex"
)

const RendererID = "superindex-carrier"
const VersionV1 = "v1"

type Renderer struct{}

func (Renderer) ID() string      { return RendererID }
func (Renderer) Version() string { return VersionV1 }

func (Renderer) Render(request renderer.Request) (renderer.Output, error) {
	if request.Width < 320 || request.Height < 320 {
		return renderer.Output{}, fmt.Errorf("superindex-carrier: dimensions must be >=320")
	}
	if request.CanonicalDocument.SourceSHA256 == "" || request.SuperIndex.SourceSHA256 == "" {
		return renderer.Output{}, fmt.Errorf("superindex-carrier: canonical document and SuperIndex are required")
	}
	if request.CanonicalDocument.SourceSHA256 != request.SourceSHA || request.SuperIndex.SourceSHA256 != request.SourceSHA {
		return renderer.Output{}, fmt.Errorf("superindex-carrier: source hash mismatch")
	}
	manifest, err := capacity.ReferencedManifest(request.CanonicalDocument, request.SuperIndex)
	if err != nil {
		return renderer.Output{}, err
	}
	const cellPixels = 36
	const marginPixels = 4
	columns := max(1, (request.Width-2*marginPixels)/cellPixels)
	glyphPNG, _, err := glyphcalc.RenderGlyphGrid(glyphcalc.GlyphsFromBytes(manifest), columns, cellPixels, marginPixels)
	if err != nil {
		return renderer.Output{}, err
	}
	glyphImage, err := renderer.DecodePNG(glyphPNG)
	if err != nil {
		return renderer.Output{}, err
	}
	if glyphImage.Bounds().Dx() > request.Width || glyphImage.Bounds().Dy() > request.Height {
		return renderer.Output{}, fmt.Errorf("superindex-carrier: referenced manifest exceeds surface")
	}
	canvas := image.NewRGBA(image.Rect(0, 0, request.Width, request.Height))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	origin := image.Pt((request.Width-glyphImage.Bounds().Dx())/2, (request.Height-glyphImage.Bounds().Dy())/2)
	draw.Draw(canvas, glyphImage.Bounds().Add(origin), glyphImage, glyphImage.Bounds().Min, draw.Src)
	pngBytes, err := renderer.EncodePNG(canvas)
	if err != nil {
		return renderer.Output{}, err
	}
	indexHash, err := superindex.Hash(request.SuperIndex)
	if err != nil {
		return renderer.Output{}, err
	}
	result := renderer.Result{Schema: renderer.ResultSchemaV1, RendererID: RendererID, Version: VersionV1, Width: request.Width, Height: request.Height, PixelSHA256: renderer.PixelHash(canvas), PNGSHA256: renderer.PNGHash(pngBytes), PNGBytes: len(pngBytes), SourceSHA256: request.SourceSHA, LayoutSeed: request.LayoutSeed, ImagePath: "images/original.png", CoverageStatus: "KNOWN", RepresentationScope: "REFERENCED", ExternalPayloadSHA256: indexHash}
	return renderer.Output{PNG: pngBytes, Image: canvas, Result: result}, nil
}
