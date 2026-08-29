package superindex

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
)

func TestBuildPopulatesFiveAxes(t *testing.T) {
	document := canonicaldocument.Document{SourceSHA256: "source", Headings: []canonicaldocument.Heading{{ID: "part", Page: 1, Text: "Part I", Level: 1}, {ID: "chapter", Page: 2, Text: "1 Cognition", Level: 2}}, Entities: []canonicaldocument.Entity{{ID: "ada", Pages: []int{2}}}, Relations: []canonicaldocument.Relation{{ID: "wrote", Kind: "author_of"}}}
	index := Build(document, map[int]string{2: "Swarm cognition"})
	if len(index.Pages) == 0 || len(index.Entities) == 0 || len(index.Relations) == 0 || len(index.Hierarchy) != 2 || len(index.Terms) != 2 {
		t.Fatalf("incomplete superindex: %+v", index)
	}
	if index.Hierarchy[1].Parent != "part" {
		t.Fatalf("chapter parent=%q", index.Hierarchy[1].Parent)
	}
}
