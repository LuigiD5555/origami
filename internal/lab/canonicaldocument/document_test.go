package canonicaldocument

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/source"
)

func TestBuildExtractsAddressableStructure(t *testing.T) {
	text := "Contents\nPart I Foundations\n1 Swarm Cognition ..... 3\nAda Lovelace and Grace Hopper\nReferences\n"
	document := Build([]byte(text), []source.PageEntry{{Page: 7, Text: text, ByteLength: len(text)}})
	if len(document.Pages) != 1 || len(document.Headings) < 3 || len(document.Tables) != 1 || len(document.References) != 1 {
		t.Fatalf("missing canonical structure: %+v", document)
	}
	if len(document.Entities) != 2 || len(document.Relations) != 2 {
		t.Fatalf("missing author relationships: entities=%v relations=%v", document.Entities, document.Relations)
	}
}
