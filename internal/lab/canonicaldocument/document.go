package canonicaldocument

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/source"
)

const SchemaV1 = "ohf.canonical-document.v1"

type Document struct {
	Schema       string      `json:"schema"`
	SourceSHA256 string      `json:"source_sha256"`
	Pages        []Page      `json:"pages"`
	Headings     []Heading   `json:"headings"`
	Tables       []Table     `json:"tables"`
	References   []Reference `json:"references"`
	Entities     []Entity    `json:"entities"`
	Relations    []Relation  `json:"relations"`
}

type Page struct {
	Number     int `json:"number"`
	ByteOffset int `json:"byte_offset"`
	ByteLength int `json:"byte_length"`
}
type Heading struct {
	ID    string `json:"id"`
	Page  int    `json:"page"`
	Text  string `json:"text"`
	Level int    `json:"level"`
}
type Table struct {
	ID   string   `json:"id"`
	Page int      `json:"page"`
	Kind string   `json:"kind"`
	Rows []string `json:"rows"`
}
type Reference struct {
	ID   string `json:"id"`
	Page int    `json:"page"`
	Text string `json:"text"`
}
type Entity struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Pages []int  `json:"pages"`
}
type Relation struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	From string `json:"from"`
	To   string `json:"to"`
	Page int    `json:"page"`
}

var chapterLine = regexp.MustCompile(`^([0-9]{1,3})\s+(.+?)(?:\s+\.{2,}\s*[0-9]+)?$`)
var email = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)

func Build(material []byte, pageEntries []source.PageEntry) Document {
	sum := sha256.Sum256(material)
	document := Document{Schema: SchemaV1, SourceSHA256: hex.EncodeToString(sum[:])}
	entityPages := map[string]map[int]bool{}
	entityKind := map[string]string{}
	for _, pageEntry := range pageEntries {
		document.Pages = append(document.Pages, Page{Number: pageEntry.Page, ByteOffset: pageEntry.ByteOffset, ByteLength: pageEntry.ByteLength})
		lines := nonemptyLines(pageEntry.Text)
		var tableRows []string
		for lineIndex, line := range lines {
			if level := headingLevel(line); level > 0 {
				document.Headings = append(document.Headings, Heading{ID: stableID("heading", pageEntry.Page, line), Page: pageEntry.Page, Text: line, Level: level})
			}
			if strings.Contains(line, "...") {
				tableRows = append(tableRows, line)
			}
			if strings.EqualFold(line, "References") || strings.HasPrefix(line, "References.") {
				document.References = append(document.References, Reference{ID: stableID("reference", pageEntry.Page, line), Page: pageEntry.Page, Text: line})
			}
			for _, address := range email.FindAllString(line, -1) {
				addEntity(entityPages, entityKind, strings.ToLower(address), "email", pageEntry.Page)
			}
			if matches := chapterLine.FindStringSubmatch(line); matches != nil {
				chapterID := stableID("chapter", pageEntry.Page, matches[1]+":"+matches[2])
				if lineIndex+1 < len(lines) {
					for _, name := range splitNames(lines[lineIndex+1]) {
						addEntity(entityPages, entityKind, name, "person", pageEntry.Page)
						document.Relations = append(document.Relations, Relation{ID: stableID("relation", pageEntry.Page, name+chapterID), Kind: "author_of", From: entityID(name), To: chapterID, Page: pageEntry.Page})
					}
				}
			}
		}
		if len(tableRows) > 0 {
			document.Tables = append(document.Tables, Table{ID: stableID("table", pageEntry.Page, strings.Join(tableRows, "\n")), Page: pageEntry.Page, Kind: "contents", Rows: tableRows})
		}
	}
	var names []string
	for name := range entityPages {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		var pages []int
		for page := range entityPages[name] {
			pages = append(pages, page)
		}
		sort.Ints(pages)
		document.Entities = append(document.Entities, Entity{ID: entityID(name), Name: name, Kind: entityKind[name], Pages: pages})
	}
	return document
}

func nonemptyLines(text string) []string {
	var lines []string
	for _, raw := range strings.Split(text, "\n") {
		if line := strings.TrimSpace(raw); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func headingLevel(line string) int {
	if strings.HasPrefix(line, "Part ") {
		return 1
	}
	if chapterLine.MatchString(line) {
		return 2
	}
	if line == "Contents" || strings.HasSuffix(line, "Contributors") || strings.HasPrefix(line, "About the ") || line == "References" {
		return 1
	}
	return 0
}

func splitNames(line string) []string {
	if len(line) > 180 || strings.Contains(line, "...") {
		return nil
	}
	normalized := strings.ReplaceAll(line, " and ", ",")
	var names []string
	for _, part := range strings.Split(normalized, ",") {
		name := strings.TrimSpace(part)
		words := strings.Fields(name)
		if len(words) >= 2 && len(words) <= 6 {
			names = append(names, name)
		}
	}
	return names
}

func addEntity(pages map[string]map[int]bool, kinds map[string]string, name, kind string, page int) {
	if pages[name] == nil {
		pages[name] = map[int]bool{}
	}
	pages[name][page] = true
	kinds[name] = kind
}

func entityID(name string) string { return stableID("entity", 0, strings.ToLower(name)) }
func stableID(kind string, page int, value string) string {
	sum := sha256.Sum256([]byte(kind + "\x00" + strconv.Itoa(page) + "\x00" + value))
	return kind + ":" + hex.EncodeToString(sum[:6])
}
