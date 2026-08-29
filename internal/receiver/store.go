package receiver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type StoredCandidate struct {
	Directory    string `json:"directory"`
	PromptPath   string `json:"prompt_path"`
	CandidatePath string `json:"candidate_path"`
	ManifestPath string `json:"manifest_path"`
}

// StoreImportedCandidate persists an imported Tlaloc proposal under Origami's
// candidate registry. It always stores the artifact as REFERENCE_CANDIDATE;
// promotion requires later Origami validation/evidence and is intentionally a
// separate operation.
func StoreImportedCandidate(root string, candidate ImportedCandidate) (StoredCandidate, error) {
	if root == "" {
		return StoredCandidate{}, fmt.Errorf("storage root is required")
	}
	if candidate.CandidateID == "" || candidate.UniversalPrompt == "" || len(candidate.Program) == 0 {
		return StoredCandidate{}, fmt.Errorf("candidate id, universal prompt and program are required")
	}
	if strings.Contains(candidate.CandidateID, "..") || strings.ContainsAny(candidate.CandidateID, `/\\`) {
		return StoredCandidate{}, fmt.Errorf("candidate id must be path-safe")
	}

	dir := filepath.Join(root, candidate.CandidateID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return StoredCandidate{}, err
	}
	promptPath := filepath.Join(dir, "MASTER_PROMPT.md")
	candidatePath := filepath.Join(dir, "candidate.json")
	manifestPath := filepath.Join(dir, "manifest.json")

	if err := os.WriteFile(promptPath, []byte(candidate.UniversalPrompt), 0644); err != nil {
		return StoredCandidate{}, err
	}
	candidateBytes, err := json.MarshalIndent(candidate, "", "  ")
	if err != nil {
		return StoredCandidate{}, err
	}
	candidateBytes = append(candidateBytes, '\n')
	if err := os.WriteFile(candidatePath, candidateBytes, 0644); err != nil {
		return StoredCandidate{}, err
	}

	sum := sha256.Sum256([]byte(candidate.UniversalPrompt))
	promptHash := hex.EncodeToString(sum[:])
	manifest := ArtifactManifest{
		Schema:             ArtifactSchemaR0,
		Status:             ArtifactReferenceCandidate,
		ReceiverContract:   ContractID,
		DistillerContract:  TlalocHybridArtifactSchema,
		CandidateID:        candidate.CandidateID,
		PromptPath:         filepath.ToSlash(promptPath),
		PromptSHA256:       promptHash,
		SourceTraceSHA256:  candidate.SourceTraceSHA256,
		WorkingWindowToken: candidate.WorkingWindow,
		Notes:              "Imported from Tlaloc; not promoted. Carrier-local BOOT/ROSETTA bindings and model evidence remain required.",
	}
	if err := ValidateArtifactManifest(manifest); err != nil {
		return StoredCandidate{}, err
	}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return StoredCandidate{}, err
	}
	manifestBytes = append(manifestBytes, '\n')
	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		return StoredCandidate{}, err
	}

	return StoredCandidate{
		Directory:     dir,
		PromptPath:    promptPath,
		CandidatePath: candidatePath,
		ManifestPath:  manifestPath,
	}, nil
}
