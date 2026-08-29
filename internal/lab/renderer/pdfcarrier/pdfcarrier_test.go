package pdfcarrier

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"github.com/LuigiD5555/origami/internal/lab/source"
)

func TestRenderProducesDeterministicReopenableCarrier(t *testing.T) {
	request := renderer.Request{Width: 640, Height: 480, Source: []byte("page one\npage two\n"), SourceSHA: "source-hash", LayoutSeed: 7, PageIndex: []source.PageEntry{{Page: 1, Text: "page one\n", ByteOffset: 0, ByteLength: 9}, {Page: 2, Text: "page two\n", ByteOffset: 9, ByteLength: 9}}}
	first, err := (Renderer{}).Render(request)
	if err != nil {
		t.Fatal(err)
	}
	second, err := (Renderer{}).Render(request)
	if err != nil {
		t.Fatal(err)
	}
	if first.Result.PNGSHA256 != second.Result.PNGSHA256 {
		t.Fatal("renderer output is not deterministic")
	}
	reopened, err := renderer.DecodePNG(first.PNG)
	if err != nil {
		t.Fatal(err)
	}
	if renderer.PixelHash(reopened) != first.Result.PixelSHA256 {
		t.Fatal("reopened pixel hash mismatch")
	}
}

func TestRenderReportsTruncatedCoverage(t *testing.T) {
	text := strings.Repeat("long content ", 2000)
	request := renderer.Request{Width: 320, Height: 320, Source: []byte(text), SourceSHA: "source-hash", PageIndex: []source.PageEntry{{Page: 1, Text: text, ByteLength: len(text)}}}
	output, err := (Renderer{}).Render(request)
	if err != nil {
		t.Fatal(err)
	}
	if output.Result.CoverageStatus != "UNKNOWN" || output.Result.TruncatedTiles != 1 || output.Result.LostBytes <= 0 {
		t.Fatalf("missing truncation evidence: %+v", output.Result)
	}
}
