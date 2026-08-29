package receiver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestR0RegistryMatchesStoredPrompt(t *testing.T) {
	root := filepath.Join("..", "..")
	manifestPath := filepath.Join(root, "receiver", "registry", "r0", "manifest.json")
	b, err := os.ReadFile(manifestPath)
	if err != nil { t.Fatal(err) }
	var manifest ArtifactManifest
	if err := json.Unmarshal(b, &manifest); err != nil { t.Fatal(err) }
	if err := ValidateArtifactManifest(manifest); err != nil { t.Fatal(err) }

	promptPath := filepath.Join(root, filepath.FromSlash(manifest.PromptPath))
	prompt, err := os.ReadFile(promptPath)
	if err != nil { t.Fatal(err) }
	sum := sha256.Sum256(prompt)
	got := hex.EncodeToString(sum[:])
	if got != manifest.PromptSHA256 {
		t.Fatalf("receiver prompt registry drift: manifest=%s actual=%s", manifest.PromptSHA256, got)
	}
	if manifest.Status == ArtifactPromoted {
		t.Fatal("R0 bootstrap must remain REFERENCE_CANDIDATE until model-facing Hybrid gates pass")
	}
}
