package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const ManifestSchemaV1 = "ohf.artifact-manifest.v1"
const ManifestFilename = "hashes.json"

type Entry struct {
	Path      string `json:"path"`
	Bytes     int64  `json:"bytes"`
	SHA256    string `json:"sha256"`
	MediaType string `json:"media_type"`
	Role      string `json:"role"`
}

type Manifest struct {
	Schema       string  `json:"schema"`
	SelfExcluded bool    `json:"self_excluded"`
	Entries      []Entry `json:"entries"`
}

// Build walks a run directory deterministically and records all regular files
// except hashes.json itself. Entries are sorted by slash-normalized relative path.
func Build(runDir string) (Manifest, error) {
	var entries []Entry
	err := filepath.WalkDir(runDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(runDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == ManifestFilename || strings.HasSuffix(rel, ".tmp") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		sha, err := FileSHA(path)
		if err != nil {
			return err
		}
		entries = append(entries, Entry{
			Path: rel, Bytes: info.Size(), SHA256: sha,
			MediaType: mediaType(rel), Role: role(rel),
		})
		return nil
	})
	if err != nil {
		return Manifest{}, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return Manifest{Schema: ManifestSchemaV1, SelfExcluded: true, Entries: entries}, nil
}

func Write(runDir string) (Manifest, error) {
	m, err := Build(runDir)
	if err != nil {
		return Manifest{}, err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return Manifest{}, err
	}
	b = append(b, '\n')
	if err := WriteExact(filepath.Join(runDir, ManifestFilename), b); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func Load(runDir string) (Manifest, error) {
	b, err := os.ReadFile(filepath.Join(runDir, ManifestFilename))
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return Manifest{}, err
	}
	if m.Schema != ManifestSchemaV1 {
		return Manifest{}, fmt.Errorf("artifact manifest: unsupported schema %q", m.Schema)
	}
	return m, nil
}

func Verify(runDir string) error {
	expected, err := Load(runDir)
	if err != nil {
		return err
	}
	actual, err := Build(runDir)
	if err != nil {
		return err
	}
	if len(expected.Entries) != len(actual.Entries) {
		return fmt.Errorf("artifact manifest: entry count mismatch got=%d want=%d", len(actual.Entries), len(expected.Entries))
	}
	for i := range expected.Entries {
		a, e := actual.Entries[i], expected.Entries[i]
		// Cryptographic integrity is defined by identity and bytes. MediaType and
		// Role are descriptive metadata and may evolve without making preserved
		// evidence appear tampered.
		if a.Path != e.Path || a.Bytes != e.Bytes || a.SHA256 != e.SHA256 {
			return fmt.Errorf("artifact manifest: mismatch for %q", e.Path)
		}
	}
	return nil
}

func FileSHA(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func WriteExact(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func mediaType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "application/json"
	case ".jsonl":
		return "application/x-ndjson"
	case ".png":
		return "image/png"
	case ".txt", ".md":
		return "text/plain"
	case ".bin":
		return "application/octet-stream"
	default:
		return "application/octet-stream"
	}
}

func role(path string) string {
	switch {
	case path == "experiment.json":
		return "experiment-spec"
	case path == "environment.json":
		return "environment"
	case path == "seeds.json":
		return "seed-manifest"
	case path == "provenance.json":
		return "provenance"
	case path == "ground_truth.json":
		return "ground-truth"
	case path == "page_index.json":
		return "address-index"
	case path == "verdict.json":
		return "verdict"
	case path == "render.json":
		return "render-metadata"
	case path == "events.jsonl":
		return "event-journal"
	case strings.HasPrefix(path, "source/"):
		return "source"
	case strings.HasPrefix(path, "images/"):
		return "carrier"
	case strings.Contains(path, "/blind/MASTER_PROMPT.md"):
		return "native-master-prompt"
	case strings.Contains(path, "/blind/carrier.png"):
		return "native-carrier-copy"
	case strings.Contains(path, "/blind/query.txt"):
		return "native-query"
	case strings.HasSuffix(path, "/bundle.json"):
		return "native-bundle-metadata"
	case strings.HasPrefix(path, "native/observations/"):
		return "native-observation"
	case strings.HasPrefix(path, "scores/"):
		return "score"
	case path == "failure.json":
		return "failure-record"
	default:
		return "artifact"
	}
}
