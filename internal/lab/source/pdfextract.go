package source

import (
	"fmt"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/spec"
	"github.com/tsawler/tabula"
)

const GeneratorPDFExtractV1 = "pdf-extract-v1"
const PageIndexSchemaV1 = "ohf.page-index.v1"

type PageEntry struct {
	Page       int    `json:"page"`
	Text       string `json:"text"`
	ByteOffset int    `json:"byte_offset"`
	ByteLength int    `json:"byte_length"`
}

type PageIndex struct {
	Schema string      `json:"schema"`
	Pages  []PageEntry `json:"pages"`
}

func GeneratePDF(experiment spec.ExperimentSpec) (Material, PageIndex, error) {
	metadata := tabula.Open(experiment.Source.Path)
	pageCount, err := metadata.PageCount()
	metadata.Close()
	if err != nil {
		return Material{}, PageIndex{}, fmt.Errorf("pdf source: open %q: %w", experiment.Source.Path, err)
	}

	var canonical strings.Builder
	index := PageIndex{Schema: PageIndexSchemaV1}
	pageStart := max(1, experiment.Source.PageStart)
	pageEnd := pageCount
	if experiment.Source.PageEnd > 0 {
		pageEnd = min(pageEnd, experiment.Source.PageEnd)
	}
	if pageStart > pageEnd {
		return Material{}, PageIndex{}, fmt.Errorf("pdf source: page range starts after document end (%d pages)", pageCount)
	}
	for pageNumber := pageStart; pageNumber <= pageEnd; pageNumber++ {
		text, _, err := tabula.Open(experiment.Source.Path).Pages(pageNumber).Text()
		if err != nil {
			return Material{}, PageIndex{}, fmt.Errorf("pdf source: extract page %d: %w", pageNumber, err)
		}
		text = normalizePageText(text)
		offset := canonical.Len()
		canonical.WriteString(text)
		index.Pages = append(index.Pages, PageEntry{Page: pageNumber, Text: text, ByteOffset: offset, ByteLength: len(text)})
	}
	if canonical.Len() == 0 {
		return Material{}, PageIndex{}, fmt.Errorf("pdf source: no extractable text")
	}
	material := Material{Bytes: []byte(canonical.String()), Generator: GeneratorPDFExtractV1, PageIndex: index.Pages}
	return material, index, nil
}

func normalizePageText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.TrimSpace(text) + "\n"
}

func VerifyPageIndex(material []byte, index PageIndex) error {
	if index.Schema != PageIndexSchemaV1 {
		return fmt.Errorf("page index: unsupported schema %q", index.Schema)
	}
	expectedOffset := 0
	for _, page := range index.Pages {
		if page.Page <= 0 || page.ByteOffset != expectedOffset || page.ByteLength != len(page.Text) {
			return fmt.Errorf("page index: invalid page %d metadata", page.Page)
		}
		end := page.ByteOffset + page.ByteLength
		if end > len(material) || string(material[page.ByteOffset:end]) != page.Text {
			return fmt.Errorf("page index: page %d does not match canonical material", page.Page)
		}
		expectedOffset = end
	}
	if expectedOffset != len(material) {
		return fmt.Errorf("page index: covered %d of %d canonical bytes", expectedOffset, len(material))
	}
	return nil
}
