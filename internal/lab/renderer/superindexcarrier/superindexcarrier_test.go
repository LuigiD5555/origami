package superindexcarrier

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"github.com/LuigiD5555/origami/internal/lab/superindex"
)

func TestRenderDeclaresReferencedPayload(t *testing.T) {
	sha := strings.Repeat("a", 64)
	request := renderer.Request{Width: 1600, Height: 2000, SourceSHA: sha, CanonicalDocument: canonicaldocument.Document{SourceSHA256: sha, Pages: []canonicaldocument.Page{{Number: 1}}}, SuperIndex: superindex.Index{SourceSHA256: sha, Pages: map[int][]string{1: {"heading"}}, Entities: map[string][]int{}, Relations: map[string][]string{}, Terms: map[string][]int{}}}
	output, err := (Renderer{}).Render(request)
	if err != nil {
		t.Fatal(err)
	}
	if output.Result.RepresentationScope != "REFERENCED" || output.Result.ExternalPayloadSHA256 == "" || output.Result.PNGBytes > 500000 {
		t.Fatalf("incorrect referenced result: %+v", output.Result)
	}
}
