package failure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/scorer"
)

const SchemaV1 = "ohf.failure.v1"

type Record struct {
	Schema         string   `json:"schema"`
	RegressionID   string   `json:"regression_id"`
	RunID          string   `json:"run_id"`
	ScoreID        string   `json:"score_id"`
	Classification string   `json:"classification"`
	SourceFiles    []string `json:"source_files"`
}

func Register(root, runID, scoreID, classification string) (Record, error) {
	runDir := filepath.Join(root, "runs", runID)
	scorePath := filepath.Join(runDir, "scores", scoreID+".json")
	b, err := os.ReadFile(scorePath)
	if err != nil { return Record{}, err }
	var s scorer.Score
	if err := json.Unmarshal(b, &s); err != nil { return Record{}, err }
	if s.Verdict != "FAIL" { return Record{}, fmt.Errorf("cannot register passing score") }
	if classification == "" { classification = classify(s) }
	obsPath, err := findObservation(runDir, s.ObservationID)
	if err != nil { return Record{}, err }
	bundleID, err := observationBundleID(obsPath)
	if err != nil { return Record{}, err }
	inputs := []string{"experiment.json", "ground_truth.json", "render.json", "images/original.png", filepath.ToSlash(strings.TrimPrefix(obsPath, runDir+string(os.PathSeparator))), filepath.ToSlash(strings.TrimPrefix(scorePath, runDir+string(os.PathSeparator))), filepath.ToSlash(filepath.Join("native", "bundles", bundleID, "bundle.json")), filepath.ToSlash(filepath.Join("native", "bundles", bundleID, "blind", "MASTER_PROMPT.md")), filepath.ToSlash(filepath.Join("native", "bundles", bundleID, "blind", "carrier.png")), filepath.ToSlash(filepath.Join("native", "bundles", bundleID, "blind", "query.txt"))}
	sort.Strings(inputs)
	regID := regressionID(runID, scoreID, classification, b)
	regDir := filepath.Join(root, "regressions", regID)
	if _, err := os.Stat(regDir); err == nil {
		if err := artifact.Verify(regDir); err != nil { return Record{}, fmt.Errorf("existing regression %s failed integrity: %w", regID, err) }
		rb, err := os.ReadFile(filepath.Join(regDir, "failure.json")); if err != nil { return Record{}, err }
		var existing Record
		if err := json.Unmarshal(rb, &existing); err != nil { return Record{}, err }
		if existing.RunID != runID || existing.ScoreID != scoreID || existing.Classification != classification { return Record{}, fmt.Errorf("existing regression %s identity collision", regID) }
		return existing, nil
	}
	for _, rel := range inputs {
		src := filepath.Join(runDir, filepath.FromSlash(rel)); data, err := os.ReadFile(src)
		if err != nil { return Record{}, fmt.Errorf("copy regression input %s: %w", rel, err) }
		if err := artifact.WriteExact(filepath.Join(regDir, "evidence", filepath.FromSlash(rel)), data); err != nil { return Record{}, err }
	}
	r := Record{Schema: SchemaV1, RegressionID: regID, RunID: runID, ScoreID: scoreID, Classification: classification, SourceFiles: inputs}
	rb, _ := json.MarshalIndent(r, "", "  ")
	if err := artifact.WriteExact(filepath.Join(regDir, "failure.json"), append(rb, '\n')); err != nil { return Record{}, err }
	if _, err := artifact.Write(regDir); err != nil { return Record{}, err }
	if err := artifact.Verify(regDir); err != nil { return Record{}, err }
	return r, nil
}

func Verify(root, regID string) error { return artifact.Verify(filepath.Join(root, "regressions", regID)) }
func classify(s scorer.Score) string { if s.FailureKind == "WRONG_OBSERVATION" || strings.HasPrefix(s.FailureKind, "EPISTEMIC_") { return "PERCEPTION" }; return "UNKNOWN" }
func regressionID(runID, scoreID, class string, score []byte) string { h := sha256.New(); io.WriteString(h, "ohf.regression.v1\x00"); io.WriteString(h, runID); h.Write([]byte{0}); io.WriteString(h, scoreID); h.Write([]byte{0}); io.WriteString(h, class); h.Write([]byte{0}); h.Write(score); return "REG-" + hex.EncodeToString(h.Sum(nil))[:12] }
func findObservation(runDir, id string) (string, error) { base := filepath.Join(runDir, "native", "observations"); var found string; err := filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error { if err != nil { return err }; if d.IsDir() { return nil }; if strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)) == id { found = p }; return nil }); if err != nil { return "", err }; if found == "" { return "", fmt.Errorf("observation %s not found", id) }; return found, nil }
func observationBundleID(path string) (string, error) { b, err := os.ReadFile(path); if err != nil { return "", err }; var x struct { BundleID string `json:"bundle_id"` }; if err := json.Unmarshal(b, &x); err != nil { return "", err }; if x.BundleID == "" { return "", fmt.Errorf("observation missing bundle_id") }; return x.BundleID, nil }
