package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	stddraw "image/draw"
	"image/png"
	"sort"
	"strings"
)

const CandidateSpecSchema = "origami.experimental-candidate.r0"
const CandidateBuildReportSchema = "origami.experimental-candidate.r0.build-report"

type CandidateMutation struct {
	Kind         string `json:"kind"`
	Target       string `json:"target"`
	Value        string `json:"value"`
	Rationale    string `json:"rationale,omitempty"`
	Experimental bool   `json:"experimental"`
}

type CandidateSpec struct {
	Schema       string              `json:"schema"`
	ID           string              `json:"id"`
	ParentSHA256 string              `json:"parent_sha256,omitempty"`
	Mutations    []CandidateMutation `json:"mutations"`
}

type CandidateBuildReport struct {
	Schema                 string              `json:"schema"`
	CandidateID            string              `json:"candidate_id"`
	ParentProfile          string              `json:"parent_profile"`
	ParentSHA256           string              `json:"parent_sha256"`
	CandidateSHA256        string              `json:"candidate_sha256"`
	ParentProgramSHA256    string              `json:"parent_program_sha256"`
	CandidateProgramSHA256 string              `json:"candidate_program_sha256"`
	ExactProgramPreserved  bool                `json:"exact_program_preserved"`
	PNGBytes               int                 `json:"png_bytes"`
	Width                  int                 `json:"width"`
	Height                 int                 `json:"height"`
	AppliedMutations       []CandidateMutation `json:"applied_mutations"`
}

type CandidateBuilderCapabilities struct {
	Schema             string   `json:"schema"`
	ParentProfiles     []string `json:"parent_profiles"`
	SupportedKinds     []string `json:"supported_kinds"`
	UnsupportedKinds   []string `json:"unsupported_kinds"`
	ExactPlaneMutation bool     `json:"exact_plane_mutation"`
	MaxMutations       int      `json:"max_mutations"`
}

func TemporalCandidateCapabilities() CandidateBuilderCapabilities {
	return CandidateBuilderCapabilities{
		Schema: "origami.experimental-candidate.r0.capabilities",
		ParentProfiles: []string{TemporalCarrierProfileID},
		SupportedKinds: []string{"CHANNEL_ROLE", "LAYOUT", "NUMERIC_STRUCTURE", "PRIMITIVE", "PROMPT", "REDUNDANCY", "TEMPORAL_STRUCTURE"},
		UnsupportedKinds: []string{"COLOR_USAGE", "DEPTH_STRUCTURE", "EMERGENT_STRUCTURE", "INTERFERENCE_STRUCTURE"},
		ExactPlaneMutation: false,
		MaxMutations: 8,
	}
}

func BuildTemporalCandidate(parentPNG []byte, spec CandidateSpec) ([]byte, CandidateBuildReport, error) {
	if err := validateCandidateSpec(spec); err != nil {
		return nil, CandidateBuildReport{}, err
	}
	parentSum := sha256.Sum256(parentPNG)
	parentSHA := hex.EncodeToString(parentSum[:])
	if spec.ParentSHA256 != "" && !strings.EqualFold(spec.ParentSHA256, parentSHA) {
		return nil, CandidateBuildReport{}, fmt.Errorf("parent sha256 mismatch: declared=%s actual=%s", spec.ParentSHA256, parentSHA)
	}
	parentDecoded, err := DecodeTemporalCarrierPNG(parentPNG)
	if err != nil {
		return nil, CandidateBuildReport{}, fmt.Errorf("parent is not a valid temporal carrier: %w", err)
	}
	if parentDecoded.Profile != TemporalCarrierProfileID {
		return nil, CandidateBuildReport{}, fmt.Errorf("unsupported parent profile %q", parentDecoded.Profile)
	}
	img, err := png.Decode(bytes.NewReader(parentPNG))
	if err != nil {
		return nil, CandidateBuildReport{}, err
	}
	gray := image.NewGray(image.Rect(0, 0, Width, Height))
	stddraw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, stddraw.Src)
	mutations := append([]CandidateMutation(nil), spec.Mutations...)
	for _, mutation := range mutations {
		if err := applyTemporalCandidateMutation(gray, parentDecoded, mutation); err != nil {
			return nil, CandidateBuildReport{}, err
		}
	}
	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&raw, gray); err != nil {
		return nil, CandidateBuildReport{}, err
	}
	if raw.Len() > FixedPNGBytes-12 {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate raw PNG exceeds frozen envelope: %d > %d", raw.Len(), FixedPNGBytes-12)
	}
	candidate, err := padPNG(raw.Bytes(), FixedPNGBytes)
	if err != nil {
		return nil, CandidateBuildReport{}, err
	}
	candidateDecoded, err := DecodeTemporalCarrierPNG(candidate)
	if err != nil {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate exact-plane verification failed: %w", err)
	}
	if candidateDecoded.ProgramSHA256 != parentDecoded.ProgramSHA256 {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate mutated exact TemporalProgram: %s != %s", candidateDecoded.ProgramSHA256, parentDecoded.ProgramSHA256)
	}
	candidateSum := sha256.Sum256(candidate)
	sort.Slice(mutations, func(i, j int) bool {
		if mutations[i].Kind == mutations[j].Kind {
			if mutations[i].Target == mutations[j].Target { return mutations[i].Value < mutations[j].Value }
			return mutations[i].Target < mutations[j].Target
		}
		return mutations[i].Kind < mutations[j].Kind
	})
	report := CandidateBuildReport{
		Schema: CandidateBuildReportSchema,
		CandidateID: spec.ID,
		ParentProfile: parentDecoded.Profile,
		ParentSHA256: parentSHA,
		CandidateSHA256: hex.EncodeToString(candidateSum[:]),
		ParentProgramSHA256: parentDecoded.ProgramSHA256,
		CandidateProgramSHA256: candidateDecoded.ProgramSHA256,
		ExactProgramPreserved: true,
		PNGBytes: len(candidate), Width: Width, Height: Height,
		AppliedMutations: mutations,
	}
	return candidate, report, nil
}

func ParseCandidateMutationsJSON(raw string) ([]CandidateMutation, error) {
	var mutations []CandidateMutation
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&mutations); err != nil { return nil, err }
	return mutations, nil
}

func validateCandidateSpec(spec CandidateSpec) error {
	if spec.Schema != "" && spec.Schema != CandidateSpecSchema {
		return fmt.Errorf("candidate schema must be %q", CandidateSpecSchema)
	}
	if strings.TrimSpace(spec.ID) == "" { return fmt.Errorf("candidate id is required") }
	if len(spec.Mutations) == 0 { return fmt.Errorf("at least one mutation is required") }
	if len(spec.Mutations) > TemporalCandidateCapabilities().MaxMutations { return fmt.Errorf("too many mutations: %d", len(spec.Mutations)) }
	for i, m := range spec.Mutations {
		if strings.TrimSpace(m.Kind) == "" || strings.TrimSpace(m.Target) == "" || strings.TrimSpace(m.Value) == "" {
			return fmt.Errorf("mutation %d requires kind, target and value", i)
		}
		if !m.Experimental { return fmt.Errorf("mutation %d must remain experimental", i) }
		if !supportedTemporalCandidateKind(strings.ToUpper(m.Kind)) { return fmt.Errorf("UNSUPPORTED_MUTATION_FOR_PROFILE: %s", m.Kind) }
	}
	return nil
}

func supportedTemporalCandidateKind(kind string) bool {
	for _, k := range TemporalCandidateCapabilities().SupportedKinds { if k == kind { return true } }
	return false
}

func applyTemporalCandidateMutation(img *image.Gray, decoded TemporalCarrierDecoded, m CandidateMutation) error {
	kind := strings.ToUpper(strings.TrimSpace(m.Kind))
	target := strings.ToUpper(strings.TrimSpace(m.Target))
	value := strings.ToUpper(strings.TrimSpace(m.Value))
	switch kind {
	case "PROMPT":
		text := "READ DECLARED ROUTE FIRST"
		if target == "ROSETTA.S2.READ_SUPERINDEX" && value == "DECLARE_T2_LOCATION_BEFORE_DECODE" { text = "S2: LOCATE T2 BEFORE DECODE" }
		drawText(img, 24, 92, 1, shortLabel(text, 38), 0)
	case "LAYOUT":
		// A visible route from the BOOT/ROSETTA box to the T2 semantic plane.
		drawSimpleLine(img, 320, 110, 320, 118)
		drawArrowTip(img, 320, 110, 320, 124)
		if target == "T1_TO_T2_ENTRY_ROUTE" && value == "EXPLICIT_DIRECTIONAL_ANCHOR" { drawText(img, 332, 112, 1, "T2", 0) }
	case "REDUNDANCY":
		drawText(img, 584, 92, 1, "T2", 0)
		drawText(img, 584, 120, 1, "T2", 0)
	case "CHANNEL_ROLE":
		box(img, 18, 114, 604, 198, 0)
		drawText(img, 470, 120, 1, "SEM ROUTE", 0)
	case "PRIMITIVE":
		// Distinct four-point marker beside the active semantic target region.
		drawSimpleLine(img, 18, 120, 26, 128)
		drawSimpleLine(img, 26, 128, 18, 136)
		drawSimpleLine(img, 18, 136, 10, 128)
		drawSimpleLine(img, 10, 128, 18, 120)
	case "TEMPORAL_STRUCTURE":
		if target == "T2_SEMANTIC_TEMPORAL_SUPERGRAPH" && value == "VISIBLE_RULE_MICROGRAMMAR_R1" {
			drawTemporalRuleMicrogrammar(img, decoded)
		} else {
			emphasizeTemporalCheckpoints(img, decoded)
		}
	case "NUMERIC_STRUCTURE":
		drawText(img, 40, 376, 1, "0", 0)
		drawText(img, 176, 376, 1, "1", 0)
		drawText(img, 312, 376, 1, "2", 0)
		drawText(img, 448, 376, 1, "3", 0)
		drawText(img, 584, 376, 1, "4", 0)
	default:
		return fmt.Errorf("UNSUPPORTED_MUTATION_FOR_PROFILE: %s", kind)
	}
	return nil
}

func emphasizeTemporalCheckpoints(img *image.Gray, decoded TemporalCarrierDecoded) {
	p := decoded.Program
	steps := p.MaxSteps
	if steps < 1 { steps = 1 }
	if steps > 32 { steps = 32 }
	if p.CheckpointEvery <= 0 { return }
	const x0, x1, y = 38, 602, 354
	for s := 0; s <= steps; s++ {
		if s%p.CheckpointEvery != 0 { continue }
		x := x0 + (x1-x0)*s/steps
		box(img, x-7, y-7, 15, 15, 0)
		box(img, x-5, y-5, 11, 11, 0xff)
		box(img, x-3, y-3, 7, 7, 0)
		box(img, x-1, y-1, 3, 3, 0xff)
	}
}
