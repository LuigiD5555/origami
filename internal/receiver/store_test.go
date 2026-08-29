package receiver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreImportedCandidateIsReferenceOnly(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	candidate := ImportedCandidate{
		CandidateID:        "candidate-1",
		UniversalPrompt:    "bootstrap\n",
		BootStrategy:       []string{"find BOOT"},
		RosettaConstraints: []string{"carrier-local"},
		Program:            []MicroRule{{ID: "m0", State: "S0", Token: "T", Action: "advance", NextState: "S1"}},
		SourceTraceSHA256:  "trace",
		WorkingWindow:      4000,
	}
	stored, err := StoreImportedCandidate(filepath.Join("receiver", "candidates"), candidate)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(stored.ManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest ArtifactManifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Status != ArtifactReferenceCandidate {
		t.Fatalf("import must not auto-promote: %s", manifest.Status)
	}
	if manifest.DistillerContract != TlalocHybridArtifactSchema {
		t.Fatalf("unexpected distiller contract %q", manifest.DistillerContract)
	}
	if manifest.PromptSHA256 == "" {
		t.Fatal("stored prompt must be hashed")
	}
}
