package capacity

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"

	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/lab/superindex"
)

const SchemaV2 = "ohf.superindex-capacity-tournament.v2"

type Measurement struct {
	Schema               string      `json:"schema"`
	Preregistration      string      `json:"preregistration"`
	CanonicalBytes       int         `json:"canonical_bytes"`
	Pages                int         `json:"pages"`
	Headings             int         `json:"headings"`
	Tables               int         `json:"tables"`
	References           int         `json:"references"`
	Entities             int         `json:"entities"`
	Relations            int         `json:"relations"`
	Directions           int         `json:"directions"`
	SurfaceGlyphCapacity int         `json:"surface_glyph_capacity"`
	MaximumPNGBytes      int         `json:"maximum_png_bytes"`
	Candidates           []Candidate `json:"candidates"`
	Conclusion           string      `json:"conclusion"`
}

type Candidate struct {
	ID                  string `json:"id"`
	Scope               string `json:"scope"`
	Dependency          string `json:"dependency"`
	EncodedBytes        int    `json:"encoded_bytes"`
	ProjectedWidth      int    `json:"projected_width"`
	ProjectedHeight     int    `json:"projected_height"`
	ProjectedPNGBytes   int    `json:"projected_png_bytes"`
	FitsSurface         bool   `json:"fits_surface"`
	FitsPNGBytes        bool   `json:"fits_png_bytes"`
	PassesDeclaredScope bool   `json:"passes_declared_scope"`
	Claim               string `json:"claim"`
}

func Measure(document canonicaldocument.Document, index superindex.Index, canonicalBytes, width, height int) (Measurement, error) {
	const cellPixels = 36
	const marginPixels = 4
	const maximumPNGBytes = 500000
	columns := max(1, (width-2*marginPixels)/cellPixels)
	rows := max(0, (height-2*marginPixels)/cellPixels)
	measurement := Measurement{Schema: SchemaV2, Preregistration: "experiments/EXP-003-superindex-capacity-001.json", CanonicalBytes: canonicalBytes, Pages: len(document.Pages), Headings: len(document.Headings), Tables: len(document.Tables), References: len(document.References), Entities: len(document.Entities), Relations: len(document.Relations), Directions: len(index.Pages) + len(index.Entities) + len(document.Relations) + len(index.Hierarchy) + len(index.Terms), SurfaceGlyphCapacity: columns * rows, MaximumPNGBytes: maximumPNGBytes}
	indexBytes, err := json.Marshal(index)
	if err != nil {
		return Measurement{}, err
	}
	compressedBytes, err := deflate(indexBytes)
	if err != nil {
		return Measurement{}, err
	}
	manifestBytes, err := ReferencedManifest(document, index)
	if err != nil {
		return Measurement{}, err
	}
	candidates := []struct {
		id, scope, dependency, claim string
		payload                      []byte
	}{
		{"direct-json-glyph-v1", "SELF_CONTAINED", "NONE", "Direct JSON projection only", indexBytes},
		{"deflate-json-glyph-v1", "SELF_CONTAINED", "Fixed zlib and JSON decoders", "Reversible compressed projection under fixed decoders", compressedBytes},
		{"referenced-payload-manifest-v1", "REFERENCED", "External SuperIndex payload resolved by SHA-256", "Visual routing manifest; not self-contained compression", manifestBytes},
	}
	for _, candidate := range candidates {
		pngBytes, renderMeta, renderErr := glyphcalc.RenderGlyphGrid(glyphcalc.GlyphsFromBytes(candidate.payload), columns, cellPixels, marginPixels)
		if renderErr != nil {
			return Measurement{}, fmt.Errorf("measure %s: %w", candidate.id, renderErr)
		}
		result := Candidate{ID: candidate.id, Scope: candidate.scope, Dependency: candidate.dependency, EncodedBytes: len(candidate.payload), ProjectedWidth: renderMeta.Width, ProjectedHeight: renderMeta.Height, ProjectedPNGBytes: len(pngBytes), Claim: candidate.claim}
		result.FitsSurface = renderMeta.Width <= width && renderMeta.Height <= height
		result.FitsPNGBytes = len(pngBytes) <= maximumPNGBytes
		result.PassesDeclaredScope = result.FitsSurface && result.FitsPNGBytes
		measurement.Candidates = append(measurement.Candidates, result)
	}
	measurement.Conclusion = "FAIL_NO_DECLARED_SCOPE_FITS"
	for _, candidate := range measurement.Candidates {
		if candidate.PassesDeclaredScope {
			measurement.Conclusion = "PASS_AT_LEAST_ONE_DECLARED_SCOPE_FITS"
			break
		}
	}
	return measurement, nil
}

func ReferencedManifest(document canonicaldocument.Document, index superindex.Index) ([]byte, error) {
	indexHash, err := superindex.Hash(index)
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Schema           string `json:"schema"`
		SourceSHA256     string `json:"source_sha256"`
		SuperIndexSHA256 string `json:"superindex_sha256"`
		Pages            int    `json:"pages"`
		Entities         int    `json:"entities"`
		Relations        int    `json:"relations"`
	}{"ohf.referenced-superindex.v1", document.SourceSHA256, indexHash, len(document.Pages), len(document.Entities), len(document.Relations)})
}

func deflate(input []byte) ([]byte, error) {
	var output bytes.Buffer
	writer, err := zlib.NewWriterLevel(&output, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(input); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
