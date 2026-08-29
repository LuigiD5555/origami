package superindex

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"unicode"

	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
)

const SchemaV1 = "ohf.superindex.v1"

func Hash(index Index) (string, error) {
	encoded, err := json.Marshal(index)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

type Index struct {
	Schema       string              `json:"schema"`
	SourceSHA256 string              `json:"source_sha256"`
	Pages        map[int][]string    `json:"pages"`
	Entities     map[string][]int    `json:"entities"`
	Relations    map[string][]string `json:"relations"`
	Hierarchy    []HierarchyNode     `json:"hierarchy"`
	Terms        map[string][]int    `json:"terms"`
}
type HierarchyNode struct {
	ID     string `json:"id"`
	Parent string `json:"parent,omitempty"`
	Page   int    `json:"page"`
	Label  string `json:"label"`
}

func Build(document canonicaldocument.Document, pageText map[int]string) Index {
	index := Index{Schema: SchemaV1, SourceSHA256: document.SourceSHA256, Pages: map[int][]string{}, Entities: map[string][]int{}, Relations: map[string][]string{}, Terms: map[string][]int{}}
	for _, heading := range document.Headings {
		index.Pages[heading.Page] = append(index.Pages[heading.Page], heading.ID)
	}
	for _, table := range document.Tables {
		index.Pages[table.Page] = append(index.Pages[table.Page], table.ID)
	}
	for _, reference := range document.References {
		index.Pages[reference.Page] = append(index.Pages[reference.Page], reference.ID)
	}
	for _, entity := range document.Entities {
		index.Entities[entity.ID] = entity.Pages
	}
	for _, relation := range document.Relations {
		index.Relations[relation.Kind] = append(index.Relations[relation.Kind], relation.ID)
	}
	currentPart := ""
	for _, heading := range document.Headings {
		parent := currentPart
		if heading.Level == 1 {
			currentPart, parent = heading.ID, ""
		}
		index.Hierarchy = append(index.Hierarchy, HierarchyNode{ID: heading.ID, Parent: parent, Page: heading.Page, Label: heading.Text})
	}
	for page, text := range pageText {
		for term := range tokenize(text) {
			index.Terms[term] = appendUnique(index.Terms[term], page)
		}
	}
	for term := range index.Terms {
		sort.Ints(index.Terms[term])
	}
	return index
}

func tokenize(text string) map[string]bool {
	terms := map[string]bool{}
	for _, term := range strings.FieldsFunc(strings.ToLower(text), func(character rune) bool { return !unicode.IsLetter(character) && !unicode.IsNumber(character) }) {
		if len([]rune(term)) >= 3 {
			terms[term] = true
		}
	}
	return terms
}
func appendUnique(values []int, value int) []int {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
