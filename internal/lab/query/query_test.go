package query

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/source"
)

func TestPlanSelectsSmallDeterministicClosure(t *testing.T) {
	pages := []source.PageEntry{
		{Page: 1, Text: "An unrelated introduction.", ByteOffset: 0, ByteLength: 27},
		{Page: 2, Text: "A thinking swarm coordinates cognition.", ByteOffset: 27, ByteLength: 39},
		{Page: 3, Text: "Thinking systems are described here.", ByteOffset: 66, ByteLength: 36},
	}
	plan := Plan(pages, "What is a thinking swarm?")
	if len(plan.Addresses) != 2 {
		t.Fatalf("addresses=%v, want two matching pages", plan.Addresses)
	}
	if plan.Addresses[0].Page != 2 || plan.ClosureBytes >= plan.TotalBytes {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}

func TestUnfoldUsesIndexedRanges(t *testing.T) {
	material := []byte("first page\nsecond page\n")
	pages := []source.PageEntry{{Page: 1, ByteOffset: 0, ByteLength: 11}, {Page: 2, ByteOffset: 11, ByteLength: 12}}
	text, err := Unfold(material, pages, QueryPlan{Addresses: []Address{{Page: 2}}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "second page") || strings.Contains(text, "first page") {
		t.Fatalf("unexpected unfolded text %q", text)
	}
}
