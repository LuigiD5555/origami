package prepare

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/baseline"
	"github.com/LuigiD5555/origami/internal/lab/canonical"
	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
	"github.com/LuigiD5555/origami/internal/lab/capacity"
	"github.com/LuigiD5555/origami/internal/lab/environment"
	"github.com/LuigiD5555/origami/internal/lab/groundtruth"
	"github.com/LuigiD5555/origami/internal/lab/rendererfingerprint"
	"github.com/LuigiD5555/origami/internal/lab/runid"
	"github.com/LuigiD5555/origami/internal/lab/seed"
	"github.com/LuigiD5555/origami/internal/lab/source"
	"github.com/LuigiD5555/origami/internal/lab/spec"
	"github.com/LuigiD5555/origami/internal/lab/superindex"
	"os"
	"path/filepath"
)

const ProtocolVersion = "R3.10-LAB"

type SeedEntry struct {
	Namespace string `json:"namespace"`
	Value     uint64 `json:"value"`
}
type SeedManifest struct {
	Schema    string      `json:"schema"`
	Algorithm string      `json:"algorithm"`
	Master    uint64      `json:"master"`
	Derived   []SeedEntry `json:"derived"`
}
type Provenance struct {
	Schema                    string `json:"schema"`
	ExperimentID              string `json:"experiment_id"`
	ExperimentSpecSHA256      string `json:"experiment_spec_sha256"`
	RunID                     string `json:"run_id"`
	RunIdentitySHA256         string `json:"run_identity_sha256"`
	SourceSHA256              string `json:"source_sha256"`
	ContinuityReferenceSHA256 string `json:"continuity_reference_sha256"`
	BaselineStatus            string `json:"baseline_status"`
	BaselineSHA256            string `json:"baseline_sha256,omitempty"`
	ProtocolVersion           string `json:"protocol_version"`
}
type Gate struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}
type FoundationVerdict struct {
	Schema string `json:"schema"`
	Phase  string `json:"phase"`
	Status string `json:"status"`
	Gates  []Gate `json:"gates"`
}
type Result struct {
	RunDir            string
	ExperimentID      string
	SpecSHA256        string
	RunID             string
	RunSHA256         string
	SourceSHA256      string
	Status            string
	CanonicalDocument canonicaldocument.Document
	SuperIndex        superindex.Index
}

func Run(root string, s spec.ExperimentSpec) (Result, error) {
	var canonicalDocument canonicaldocument.Document
	var materializedSuperIndex superindex.Index
	expID, specHash, err := runid.ExperimentID(s)
	if err != nil {
		return Result{}, err
	}
	truth, material, err := groundtruth.Build(s)
	if err != nil {
		return Result{}, err
	}
	if err := groundtruth.Verify(truth, material.Bytes); err != nil {
		return Result{}, err
	}
	rf, err := rendererfingerprint.Compute(root, s.Renderer.ID, s.Renderer.Version)
	if err != nil {
		return Result{}, err
	}
	mi := runid.MaterializationIdentity{Schema: runid.MaterializationSchemaV1, ExperimentSpecSHA256: specHash, SourceSHA256: truth.Source.SHA256, RendererID: s.Renderer.ID, RendererVersion: s.Renderer.Version, RendererBinarySHA256: rf, ProtocolVersion: ProtocolVersion, BaselineVersion: s.Baseline}
	rid, rhash, err := runid.RunID(mi)
	if err != nil {
		return Result{}, err
	}
	dir := filepath.Join(root, "runs", rid)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, err
	}
	cb, err := canonical.SpecBytes(s)
	if err != nil {
		return Result{}, err
	}
	if err := artifact.WriteExact(filepath.Join(dir, "experiment.json"), append(cb, '\n')); err != nil {
		return Result{}, err
	}
	if err := writeJSON(filepath.Join(dir, "materialization.json"), mi); err != nil {
		return Result{}, err
	}
	env := environment.Capture(root)
	if err := writeJSON(filepath.Join(dir, "environment.json"), env); err != nil {
		return Result{}, err
	}
	ns := []string{source.SeedNamespace, "layout", "motifs", "negatives", "degradation", "ordering", "sampling"}
	sm := SeedManifest{Schema: "ohf.seed-manifest.v1", Algorithm: seed.Algorithm, Master: s.Seed.Value}
	for _, n := range ns {
		sm.Derived = append(sm.Derived, SeedEntry{Namespace: n, Value: seed.Derive(s.Seed.Value, n)})
	}
	if err := writeJSON(filepath.Join(dir, "seeds.json"), sm); err != nil {
		return Result{}, err
	}
	if err := writeJSON(filepath.Join(dir, "ground_truth.json"), truth); err != nil {
		return Result{}, err
	}
	if err := artifact.WriteExact(filepath.Join(dir, "source", "source.bin"), material.Bytes); err != nil {
		return Result{}, err
	}
	if len(material.PageIndex) > 0 {
		index := source.PageIndex{Schema: source.PageIndexSchemaV1, Pages: material.PageIndex}
		if err := source.VerifyPageIndex(material.Bytes, index); err != nil {
			return Result{}, err
		}
		if err := writeJSON(filepath.Join(dir, "page_index.json"), index); err != nil {
			return Result{}, err
		}
		document := canonicaldocument.Build(material.Bytes, material.PageIndex)
		canonicalDocument = document
		if document.SourceSHA256 != truth.Source.SHA256 {
			return Result{}, fmt.Errorf("canonical document source hash mismatch")
		}
		if err := writeJSON(filepath.Join(dir, "canonical_document.json"), document); err != nil {
			return Result{}, err
		}
		pageText := make(map[int]string, len(material.PageIndex))
		for _, page := range material.PageIndex {
			pageText[page.Page] = page.Text
		}
		super := superindex.Build(document, pageText)
		materializedSuperIndex = super
		if err := writeJSON(filepath.Join(dir, "superindex.json"), super); err != nil {
			return Result{}, err
		}
		measurement, err := capacity.Measure(document, super, len(material.Bytes), s.Renderer.Width, s.Renderer.Height)
		if err != nil {
			return Result{}, err
		}
		if err := writeJSON(filepath.Join(dir, "capacity.json"), measurement); err != nil {
			return Result{}, err
		}
	}
	cont := filepath.Join(root, "baseline", "r310", "reference", "ORIGAMI_OHF_CONTINUIDAD_R310.txt")
	csha, _ := fileSHA(cont)
	bv := baseline.Verify(root)
	prov := Provenance{Schema: "ohf.provenance.v2", ExperimentID: expID, ExperimentSpecSHA256: specHash, RunID: rid, RunIdentitySHA256: rhash, SourceSHA256: truth.Source.SHA256, ContinuityReferenceSHA256: csha, BaselineStatus: bv.Status, BaselineSHA256: bv.SHA256, ProtocolVersion: ProtocolVersion}
	if err := writeJSON(filepath.Join(dir, "provenance.json"), prov); err != nil {
		return Result{}, err
	}
	status := "BLOCKED_BASELINE"
	if bv.Status == "PASS" {
		status = "READY_FOR_RENDERER_FOUNDATION"
	}
	verdict := FoundationVerdict{Schema: "ohf.foundation-verdict.v2", Phase: "F1.0-F1.7", Status: status, Gates: []Gate{{ID: "G01_BASELINE_R310_PASS", Status: gateStatus(bv.Status == "PASS"), Detail: detailBaseline(bv)}, {ID: "G04_SPEC_CANONICALIZATION_PASS", Status: "PASS", Detail: specHash}, {ID: "G05_SEED_REPRODUCIBILITY_PASS", Status: "PASS", Detail: seed.Algorithm}, {ID: "G05B_CONTENT_RUN_ID_PASS", Status: "PASS", Detail: rhash}, {ID: "G06_ARTIFACT_HASHING_PASS", Status: "PENDING", Detail: artifact.ManifestFilename}, {ID: "G07_GROUND_TRUTH_FIRST_PASS", Status: "PASS", Detail: truth.Source.SHA256}, {ID: "G07B_SOURCE_TRUTH_VERIFY_PASS", Status: "PASS", Detail: fmt.Sprintf("blocks=%d", len(truth.Source.Blocks))}, {ID: "G06B_PROVENANCE_CAPTURE_PASS", Status: "PASS", Detail: env.GitCommit}}}
	if err := writeJSON(filepath.Join(dir, "verdict.json"), verdict); err != nil {
		return Result{}, err
	}
	if _, err := artifact.Write(dir); err != nil {
		return Result{}, err
	}
	if err := artifact.Verify(dir); err != nil {
		return Result{}, err
	}
	for i := range verdict.Gates {
		if verdict.Gates[i].ID == "G06_ARTIFACT_HASHING_PASS" {
			verdict.Gates[i].Status = "PASS"
			verdict.Gates[i].Detail = "manifest verified"
		}
	}
	if err := writeJSON(filepath.Join(dir, "verdict.json"), verdict); err != nil {
		return Result{}, err
	}
	if _, err := artifact.Write(dir); err != nil {
		return Result{}, err
	}
	if err := artifact.Verify(dir); err != nil {
		return Result{}, err
	}
	return Result{RunDir: dir, ExperimentID: expID, SpecSHA256: specHash, RunID: rid, RunSHA256: rhash, SourceSHA256: truth.Source.SHA256, Status: verdict.Status, CanonicalDocument: canonicalDocument, SuperIndex: materializedSuperIndex}, nil
}
func gateStatus(ok bool) string {
	if ok {
		return "PASS"
	}
	return "FAIL"
}
func detailBaseline(v baseline.Verification) string {
	if len(v.Problems) == 0 {
		return v.SHA256
	}
	return fmt.Sprint(v.Problems)
}
func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return artifact.WriteExact(path, append(b, '\n'))
}
func fileSHA(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
