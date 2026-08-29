package hybridbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

const SchemaR0 = "origami.hybrid-bundle.r0"

type FileIdentity struct {
	Path   string `json:"path"`
	Bytes  int    `json:"bytes"`
	SHA256 string `json:"sha256"`
}

type Bundle struct {
	Schema             string         `json:"schema"`
	BundleID           string         `json:"bundle_id"`
	RunID              string         `json:"run_id"`
	PublicRelativePath string         `json:"public_relative_path"`
	Files              []FileIdentity `json:"files"`
	Isolation          string         `json:"isolation"`
}

// Build creates the exact public file boundary for a Hybrid receiver trial.
// Runtime/tool implementation stays outside this directory and is authorized
// only through model_packet.json. The private source/oracle never enters it.
func Build(root, runID, promptPath, carrierPath, packetPath, queryPath string) (Bundle, error) {
	inputs := []struct {
		name string
		path string
	}{
		{"MASTER_PROMPT.md", promptPath},
		{"carrier.png", carrierPath},
		{"model_packet.json", packetPath},
		{"query.txt", queryPath},
	}
	data := make(map[string][]byte, len(inputs))
	for _, input := range inputs {
		b, err := os.ReadFile(input.path)
		if err != nil {
			return Bundle{}, fmt.Errorf("read %s: %w", input.name, err)
		}
		data[input.name] = b
	}

	id := bundleID(runID, data)
	base := filepath.Join(root, "runs", runID, "hybrid", "bundles", id)
	public := filepath.Join(base, "public")
	if err := os.RemoveAll(public); err != nil { return Bundle{}, err }
	if err := os.MkdirAll(public, 0o755); err != nil { return Bundle{}, err }

	var files []FileIdentity
	for _, input := range inputs {
		b := data[input.name]
		path := filepath.Join(public, input.name)
		if err := os.WriteFile(path, b, 0o644); err != nil { return Bundle{}, err }
		files = append(files, FileIdentity{Path: input.name, Bytes: len(b), SHA256: hash(b)})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	bundle := Bundle{
		Schema: SchemaR0,
		BundleID: id,
		RunID: runID,
		PublicRelativePath: filepath.ToSlash(filepath.Join("hybrid", "bundles", id, "public")),
		Files: files,
		Isolation: "MASTER_PROMPT_PLUS_ONE_IMAGE_PLUS_MODEL_PACKET_PLUS_QUERY_AND_DECLARED_TOOLS_ONLY",
	}
	meta, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil { return Bundle{}, err }
	meta = append(meta, '\n')
	if err := os.WriteFile(filepath.Join(base, "bundle.json"), meta, 0o644); err != nil { return Bundle{}, err }
	if err := VerifyPublic(public); err != nil { return Bundle{}, err }
	return bundle, nil
}

func VerifyPublic(public string) error {
	expected := map[string]bool{
		"MASTER_PROMPT.md": false,
		"carrier.png": false,
		"model_packet.json": false,
		"query.txt": false,
	}
	err := filepath.WalkDir(public, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return err }
		if path == public { return nil }
		rel, err := filepath.Rel(public, path)
		if err != nil { return err }
		name := filepath.ToSlash(rel)
		if d.IsDir() { return fmt.Errorf("hybrid public boundary contains unexpected directory %q", name) }
		if !d.Type().IsRegular() { return fmt.Errorf("hybrid public boundary contains non-regular entry %q", name) }
		if _, ok := expected[name]; !ok { return fmt.Errorf("hybrid public boundary contains forbidden input %q", name) }
		expected[name] = true
		return nil
	})
	if err != nil { return err }
	for name, seen := range expected {
		if !seen { return fmt.Errorf("hybrid public boundary missing %s", name) }
	}
	return nil
}

func bundleID(runID string, data map[string][]byte) string {
	h := sha256.New()
	h.Write([]byte(SchemaR0))
	h.Write([]byte{0})
	h.Write([]byte(runID))
	keys := make([]string, 0, len(data))
	for key := range data { keys = append(keys, key) }
	sort.Strings(keys)
	for _, key := range keys {
		h.Write([]byte{0})
		h.Write([]byte(key))
		h.Write([]byte{0})
		h.Write(data[key])
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func hash(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
