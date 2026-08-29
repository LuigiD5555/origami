package query

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/LuigiD5555/origami/internal/lab/source"
)

const maximumPages = 3

type Address struct {
	Page   int    `json:"page"`
	Reason string `json:"reason"`
}

type QueryPlan struct {
	Query        string    `json:"query"`
	Addresses    []Address `json:"addresses"`
	ClosureBytes int       `json:"closure_bytes"`
	TotalBytes   int       `json:"total_bytes"`
}

type candidate struct {
	entry source.PageEntry
	score int
	terms []string
}

func Plan(pageIndex []source.PageEntry, query string) QueryPlan {
	plan := QueryPlan{Query: query}
	queryTerms := tokenize(query)
	var candidates []candidate
	for _, page := range pageIndex {
		plan.TotalBytes += page.ByteLength
		pageTerms := tokenize(page.Text)
		var matched []string
		for term := range queryTerms {
			if pageTerms[term] {
				matched = append(matched, term)
			}
		}
		if len(matched) > 0 {
			sort.Strings(matched)
			candidates = append(candidates, candidate{entry: page, score: len(matched), terms: matched})
		}
	}
	sort.SliceStable(candidates, func(first, second int) bool {
		if candidates[first].score != candidates[second].score {
			return candidates[first].score > candidates[second].score
		}
		return candidates[first].entry.Page < candidates[second].entry.Page
	})
	if len(candidates) > maximumPages {
		candidates = candidates[:maximumPages]
	}
	for _, selected := range candidates {
		plan.Addresses = append(plan.Addresses, Address{Page: selected.entry.Page, Reason: "token overlap: " + strings.Join(selected.terms, ", ")})
		plan.ClosureBytes += selected.entry.ByteLength
	}
	return plan
}

func Unfold(material []byte, pageIndex []source.PageEntry, plan QueryPlan) (string, error) {
	byPage := make(map[int]source.PageEntry, len(pageIndex))
	for _, page := range pageIndex {
		byPage[page.Page] = page
	}
	var unfolded strings.Builder
	for _, address := range plan.Addresses {
		page, exists := byPage[address.Page]
		if !exists {
			return "", fmt.Errorf("query: page %d is absent from index", address.Page)
		}
		end := page.ByteOffset + page.ByteLength
		if page.ByteOffset < 0 || page.ByteLength < 0 || end > len(material) {
			return "", fmt.Errorf("query: page %d has invalid byte range", address.Page)
		}
		fmt.Fprintf(&unfolded, "--- PAGE %d ---\n", page.Page)
		unfolded.Write(material[page.ByteOffset:end])
	}
	return unfolded.String(), nil
}

func tokenize(text string) map[string]bool {
	terms := make(map[string]bool)
	for _, term := range strings.FieldsFunc(strings.ToLower(text), func(character rune) bool { return !unicode.IsLetter(character) && !unicode.IsNumber(character) }) {
		if len([]rune(term)) < 3 || stopwords[term] {
			continue
		}
		terms[term] = true
	}
	return terms
}

var stopwords = map[string]bool{
	"and": true, "are": true, "con": true, "del": true, "for": true, "las": true, "los": true,
	"para": true, "por": true, "que": true, "the": true, "una": true, "uno": true, "what": true,
	"qué": true, "como": true, "cómo": true, "this": true, "that": true,
}
