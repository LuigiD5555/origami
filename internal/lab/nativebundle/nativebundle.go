package nativebundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/evaluator"
)

const SchemaV1 = "ohf.native-bundle.v1"

type FileIdentity struct {
	Path   string `json:"path"`
	Bytes  int    `json:"bytes"`
	SHA256 string `json:"sha256"`
}
type Bundle struct {
	Schema            string         `json:"schema"`
	BundleID          string         `json:"bundle_id"`
	RunID             string         `json:"run_id"`
	BlindRelativePath string         `json:"blind_relative_path"`
	Files             []FileIdentity `json:"files"`
	Isolation         string         `json:"isolation"`
}

func Build(root, runID, promptPath, queryPath string) (Bundle, error) {
	prompt, err := os.ReadFile(promptPath)
	if err != nil {
		return Bundle{}, fmt.Errorf("read prompt: %w", err)
	}
	query, err := os.ReadFile(queryPath)
	if err != nil {
		return Bundle{}, fmt.Errorf("read query: %w", err)
	}
	return BuildWithQuery(root, runID, prompt, query)
}
func BuildWithQuery(root, runID string, prompt, query []byte) (Bundle, error) {
	if err := evaluator.ValidateIsolation(); err != nil {
		return Bundle{}, err
	}
	runDir := filepath.Join(root, "runs", runID)
	carrierPath := filepath.Join(runDir, "images", "original.png")
	carrier, err := os.ReadFile(carrierPath)
	if err != nil {
		return Bundle{}, fmt.Errorf("read carrier: %w", err)
	}
	id := bundleID(runID, prompt, carrier, query)
	base := filepath.Join(runDir, "native", "bundles", id)
	blind := filepath.Join(base, "blind")
	if err := os.RemoveAll(blind); err != nil {
		return Bundle{}, err
	}
	if err := os.MkdirAll(blind, 0o755); err != nil {
		return Bundle{}, err
	}
	files := []struct {
		name string
		data []byte
	}{{"MASTER_PROMPT.md", prompt}, {"carrier.png", carrier}, {"query.txt", query}}
	ids := make([]FileIdentity, 0, 3)
	for _, f := range files {
		p := filepath.Join(blind, f.name)
		if err := artifact.WriteExact(p, f.data); err != nil {
			return Bundle{}, err
		}
		ids = append(ids, FileIdentity{Path: f.name, Bytes: len(f.data), SHA256: hash(f.data)})
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i].Path < ids[j].Path })
	b := Bundle{Schema: SchemaV1, BundleID: id, RunID: runID, BlindRelativePath: filepath.ToSlash(filepath.Join("native", "bundles", id, "blind")), Files: ids, Isolation: "EXACTLY_MASTER_PROMPT_PLUS_ONE_IMAGE_PLUS_QUERY"}
	meta, _ := json.MarshalIndent(b, "", "  ")
	meta = append(meta, '\n')
	if err := artifact.WriteExact(filepath.Join(base, "bundle.json"), meta); err != nil {
		return Bundle{}, err
	}
	if err := VerifyBlind(blind); err != nil {
		return Bundle{}, err
	}
	if _, err := artifact.Write(runDir); err != nil {
		return Bundle{}, err
	}
	return b, nil
}
func Verify(root, runID, bundleID string) (Bundle, error) {
	base := filepath.Join(root, "runs", runID, "native", "bundles", bundleID)
	bts, err := os.ReadFile(filepath.Join(base, "bundle.json"))
	if err != nil {
		return Bundle{}, err
	}
	var b Bundle
	if err := json.Unmarshal(bts, &b); err != nil {
		return Bundle{}, err
	}
	if b.Schema != SchemaV1 || b.RunID != runID || b.BundleID != bundleID {
		return Bundle{}, fmt.Errorf("bundle metadata identity mismatch")
	}
	blind := filepath.Join(base, "blind")
	if err := VerifyBlind(blind); err != nil {
		return Bundle{}, err
	}
	for _, e := range b.Files {
		data, err := os.ReadFile(filepath.Join(blind, e.Path))
		if err != nil {
			return Bundle{}, err
		}
		if len(data) != e.Bytes || hash(data) != e.SHA256 {
			return Bundle{}, fmt.Errorf("blind file identity mismatch: %s", e.Path)
		}
	}
	return b, nil
}
func VerifyBlind(blind string) error {
	expected := map[string]bool{"MASTER_PROMPT.md": true, "carrier.png": true, "query.txt": true}
	seen := map[string]bool{}
	err := filepath.WalkDir(blind, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == blind {
			return nil
		}
		rel, err := filepath.Rel(blind, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return fmt.Errorf("blind boundary contains unexpected directory %q", filepath.ToSlash(rel))
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("blind boundary contains non-regular entry %q", filepath.ToSlash(rel))
		}
		name := filepath.ToSlash(rel)
		if !expected[name] {
			return fmt.Errorf("blind boundary contains forbidden fourth input %q", name)
		}
		seen[name] = true
		return nil
	})
	if err != nil {
		return err
	}
	if len(seen) != 3 {
		return fmt.Errorf("blind boundary expected exactly 3 files, got %d", len(seen))
	}
	for k := range expected {
		if !seen[k] {
			return fmt.Errorf("blind boundary missing %s", k)
		}
	}
	return nil
}
func bundleID(runID string, prompt, carrier, query []byte) string {
	h := sha256.New()
	h.Write([]byte("ohf.native-bundle.v1\x00"))
	h.Write([]byte(runID))
	h.Write([]byte{0})
	h.Write(prompt)
	h.Write([]byte{0})
	h.Write(carrier)
	h.Write([]byte{0})
	h.Write(query)
	return hex.EncodeToString(h.Sum(nil))[:16]
}
func hash(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
