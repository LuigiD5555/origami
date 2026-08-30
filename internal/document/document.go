package document

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/LuigiD5555/origami/internal/memory"
)

const SchemaR0 = "origami.document.r0"

type BBox struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

type Block struct {
	ID         string   `json:"id"`
	Page       int      `json:"page"`
	Order      int      `json:"order"`
	Kind       string   `json:"kind"`
	Text       string   `json:"text"`
	BBox       *BBox    `json:"bbox,omitempty"`
	SourceHash string   `json:"source_hash"`
	Relations  []string `json:"relations,omitempty"`
}

type Page struct {
	Number int     `json:"number"`
	Blocks []Block `json:"blocks"`
}

type Document struct {
	Schema     string `json:"schema"`
	SourceID   string `json:"source_id"`
	SourceHash string `json:"source_hash"`
	MediaType  string `json:"media_type"`
	Plane      string `json:"plane"`
	Pages      []Page `json:"pages"`
}

type Bundle struct {
	Schema   string        `json:"schema"`
	Document Document      `json:"document"`
	Nodes    []memory.Node `json:"nodes"`
}

func FromText(sourceID, mediaType string, data []byte) (Document, error) {
	if strings.TrimSpace(sourceID) == "" {
		return Document{}, fmt.Errorf("source_id is required")
	}
	if mediaType == "" {
		mediaType = "text/plain"
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	if strings.TrimSpace(text) == "" {
		return Document{}, fmt.Errorf("source contains no extractable text")
	}
	pageTexts := strings.Split(text, "\f")
	pages := make([]Page, 0, len(pageTexts))
	for pageIndex, pageText := range pageTexts {
		blocks := splitBlocks(pageText, pageIndex+1)
		if len(blocks) == 0 && len(pageTexts) > 1 {
			continue
		}
		pages = append(pages, Page{Number: pageIndex + 1, Blocks: blocks})
	}
	doc := Document{
		Schema:     SchemaR0,
		SourceID:   sourceID,
		SourceHash: contentHash(data),
		MediaType:  mediaType,
		Plane:      "EXACT_SOURCE_PLANE",
		Pages:      pages,
	}
	if err := Validate(doc); err != nil {
		return Document{}, err
	}
	return doc, nil
}

func LoadJSON(data []byte) (Document, error) {
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return Document{}, err
	}
	if err := Validate(doc); err != nil {
		return Document{}, err
	}
	return doc, nil
}

func Validate(doc Document) error {
	if doc.Schema != SchemaR0 {
		return fmt.Errorf("document schema must be %q", SchemaR0)
	}
	if strings.TrimSpace(doc.SourceID) == "" || strings.TrimSpace(doc.SourceHash) == "" {
		return fmt.Errorf("source_id and source_hash are required")
	}
	if doc.Plane != "EXACT_SOURCE_PLANE" {
		return fmt.Errorf("document plane must be EXACT_SOURCE_PLANE")
	}
	seen := map[string]struct{}{}
	lastPage := 0
	for _, page := range doc.Pages {
		if page.Number <= lastPage {
			return fmt.Errorf("pages must be strictly ordered")
		}
		lastPage = page.Number
		for _, block := range page.Blocks {
			if block.Page != page.Number || block.Order <= 0 || block.ID == "" {
				return fmt.Errorf("invalid block identity on page %d", page.Number)
			}
			identity := fmt.Sprintf("%d/%s", page.Number, block.ID)
			if _, ok := seen[identity]; ok {
				return fmt.Errorf("duplicate block identity %q", identity)
			}
			seen[identity] = struct{}{}
			if strings.TrimSpace(block.Text) == "" {
				return fmt.Errorf("block %q has empty text", block.ID)
			}
			if block.SourceHash == "" || block.SourceHash != contentHash([]byte(block.Text)) {
				return fmt.Errorf("block %q source hash mismatch", block.ID)
			}
		}
	}
	if len(seen) == 0 {
		return fmt.Errorf("document has no text blocks")
	}
	return nil
}

func BuildBundle(doc Document) (Bundle, error) {
	if err := Validate(doc); err != nil {
		return Bundle{}, err
	}
	nodes, err := ToMemoryNodes(doc)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{Schema: SchemaR0 + ".bundle", Document: doc, Nodes: nodes}, nil
}

func ToMemoryNodes(doc Document) ([]memory.Node, error) {
	if err := Validate(doc); err != nil {
		return nil, err
	}
	carrierID := cleanID(doc.SourceID)
	rootAddress := memory.CanonicalAddress(carrierID, "source", "root")
	pageAddresses := make([]string, 0, len(doc.Pages))
	var allText strings.Builder
	for _, page := range doc.Pages {
		for _, block := range page.Blocks {
			allText.WriteString(block.Text)
			allText.WriteByte(' ')
		}
	}
	nodes := []memory.Node{{
		Address:       rootAddress,
		CarrierID:     carrierID,
		ClusterID:     "source",
		Kind:          "source-root",
		SemanticClass: "exact-source-index",
		Label:         doc.SourceID,
		Levels: map[memory.Fidelity]memory.ContentLevel{
			memory.FidelityAbstract: {Text: keywordMetadata(allText.String(), 24)},
		},
		Provenance:      memory.Provenance{SourceID: doc.SourceID, SourceRef: doc.SourceID, SourceHash: doc.SourceHash},
		VerificationRef: doc.SourceHash,
	}}

	for _, page := range doc.Pages {
		pageID := fmt.Sprintf("page-%06d", page.Number)
		pageAddress := memory.CanonicalAddress(carrierID, "source", pageID)
		pageAddresses = append(pageAddresses, pageAddress)
		blockAddresses := make([]string, 0, len(page.Blocks))
		var pageText strings.Builder
		for _, block := range page.Blocks {
			pageText.WriteString(block.Text)
			pageText.WriteByte(' ')
			blockAddress := memory.CanonicalAddress(carrierID, pageID, block.ID)
			blockAddresses = append(blockAddresses, blockAddress)
		}
		nodes = append(nodes, memory.Node{
			Address:       pageAddress,
			CarrierID:     carrierID,
			ClusterID:     "source",
			Kind:          "page",
			SemanticClass: "source-structure",
			Label:         fmt.Sprintf("page %d", page.Number),
			Levels: map[memory.Fidelity]memory.ContentLevel{
				memory.FidelityAbstract: {Text: keywordMetadata(pageText.String(), 16)},
			},
			Parent:          rootAddress,
			Children:        blockAddresses,
			Provenance:      memory.Provenance{SourceID: doc.SourceID, SourceRef: fmt.Sprintf("%s#page=%d", doc.SourceID, page.Number), SourceHash: doc.SourceHash},
			VerificationRef: doc.SourceHash,
		})

		for i, block := range page.Blocks {
			blockAddress := blockAddresses[i]
			edges := []memory.Edge{}
			if i > 0 {
				edges = append(edges, memory.Edge{Relation: "previous", Target: blockAddresses[i-1], Weight: 1})
			}
			if i+1 < len(blockAddresses) {
				edges = append(edges, memory.Edge{Relation: "next", Target: blockAddresses[i+1], Weight: 1})
			}
			for _, relation := range block.Relations {
				if relation != "" {
					edges = append(edges, memory.Edge{Relation: "declared", Target: relation, Weight: 1})
				}
			}
			nodes = append(nodes, memory.Node{
				Address:       blockAddress,
				CID:           memory.ContentID(block.Text),
				CarrierID:     carrierID,
				ClusterID:     pageID,
				Kind:          "source-block",
				SemanticClass: block.Kind,
				Label:         fmt.Sprintf("page %d block %d %s", page.Number, block.Order, block.Kind),
				Levels: map[memory.Fidelity]memory.ContentLevel{
					memory.FidelityAbstract: {Text: keywordMetadata(block.Text, 12)},
					memory.FidelityEvidence: {Text: block.Text, SourceRef: blockAddress, Verified: true},
					memory.FidelityExact:    {Text: block.Text, SourceRef: blockAddress, Verified: true},
				},
				Parent:          pageAddress,
				Edges:           edges,
				Provenance:      memory.Provenance{SourceID: doc.SourceID, SourceRef: fmt.Sprintf("%s#page=%d&block=%d", doc.SourceID, page.Number, block.Order), SourceHash: block.SourceHash},
				VerificationRef: block.SourceHash,
			})
		}
	}
	nodes[0].Children = pageAddresses
	return nodes, nil
}

func splitBlocks(pageText string, page int) []Block {
	lines := strings.Split(pageText, "\n")
	var blocks []Block
	var current []string
	flush := func() {
		text := strings.TrimSpace(strings.Join(current, "\n"))
		current = nil
		if text == "" {
			return
		}
		order := len(blocks) + 1
		kind := "text"
		trimmed := strings.TrimSpace(text)
		if strings.HasPrefix(trimmed, "#") {
			kind = "heading"
		}
		id := fmt.Sprintf("block-%06d", order)
		blocks = append(blocks, Block{ID: id, Page: page, Order: order, Kind: kind, Text: text, SourceHash: contentHash([]byte(text))})
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		current = append(current, line)
	}
	flush()
	return blocks
}

func keywordMetadata(text string, maxTerms int) string {
	if maxTerms <= 0 {
		maxTerms = 12
	}
	counts := map[string]int{}
	for _, term := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') }) {
		term = strings.TrimSpace(term)
		if len([]rune(term)) < 3 || stopWord(term) {
			continue
		}
		counts[term]++
	}
	type pair struct { term string; count int }
	pairs := make([]pair, 0, len(counts))
	for term, count := range counts { pairs = append(pairs, pair{term, count}) }
	sort.Slice(pairs, func(i, j int) bool { if pairs[i].count == pairs[j].count { return pairs[i].term < pairs[j].term }; return pairs[i].count > pairs[j].count })
	if len(pairs) > maxTerms { pairs = pairs[:maxTerms] }
	terms := make([]string, 0, len(pairs))
	for _, item := range pairs { terms = append(terms, item.term) }
	return "keywords: " + strings.Join(terms, " ")
}

func stopWord(term string) bool {
	switch term {
	case "the", "and", "for", "with", "from", "that", "this", "into", "are", "was", "were", "una", "uno", "unos", "unas", "que", "qué", "como", "cómo", "por", "para", "con", "del", "las", "los", "esto", "esta", "este", "entre", "sobre", "desde":
		return true
	default:
		return false
	}
}

func contentHash(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func cleanID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" { return "source" }
	var b strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' { b.WriteRune(r) } else { b.WriteByte('-') }
	}
	return strings.Trim(b.String(), "-")
}
