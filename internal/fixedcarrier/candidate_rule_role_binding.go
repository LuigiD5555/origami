package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	stddraw "image/draw"
	"image/png"
	"strings"
)

const VisibleRuleRoleBindingR1 = "VISIBLE_RULE_ROLE_BINDING_R1"

// BuildRuleRoleBindingCandidate is the deterministic R5 materializer. It changes
// only the visible T2 rule-role surface and verifies that the exact TemporalProgram
// plane remains byte-decodable with the same program SHA.
func BuildRuleRoleBindingCandidate(parentPNG []byte, id string) ([]byte, CandidateBuildReport, error) {
	parentSum := sha256.Sum256(parentPNG)
	parentSHA := hex.EncodeToString(parentSum[:])
	parentDecoded, err := DecodeTemporalCarrierPNG(parentPNG)
	if err != nil { return nil, CandidateBuildReport{}, fmt.Errorf("parent is not a valid temporal carrier: %w", err) }
	img, err := png.Decode(bytes.NewReader(parentPNG))
	if err != nil { return nil, CandidateBuildReport{}, err }
	gray := image.NewGray(image.Rect(0, 0, Width, Height))
	stddraw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, stddraw.Src)
	drawVisibleRuleRoleBindingR1(gray, parentDecoded)
	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&raw, gray); err != nil { return nil, CandidateBuildReport{}, err }
	if raw.Len() > FixedPNGBytes-12 { return nil, CandidateBuildReport{}, fmt.Errorf("candidate raw PNG exceeds frozen envelope: %d > %d", raw.Len(), FixedPNGBytes-12) }
	candidate, err := padPNG(raw.Bytes(), FixedPNGBytes)
	if err != nil { return nil, CandidateBuildReport{}, err }
	candidateDecoded, err := DecodeTemporalCarrierPNG(candidate)
	if err != nil { return nil, CandidateBuildReport{}, fmt.Errorf("candidate exact-plane verification failed: %w", err) }
	if candidateDecoded.ProgramSHA256 != parentDecoded.ProgramSHA256 { return nil, CandidateBuildReport{}, fmt.Errorf("candidate mutated exact TemporalProgram") }
	candidateSum := sha256.Sum256(candidate)
	mutation := CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"RULE_ROLE_BINDING",Value:VisibleRuleRoleBindingR1,Experimental:true}
	return candidate, CandidateBuildReport{Schema:CandidateBuildReportSchema,CandidateID:id,ParentProfile:parentDecoded.Profile,ParentSHA256:parentSHA,CandidateSHA256:hex.EncodeToString(candidateSum[:]),ParentProgramSHA256:parentDecoded.ProgramSHA256,CandidateProgramSHA256:candidateDecoded.ProgramSHA256,ExactProgramPreserved:true,PNGBytes:len(candidate),Width:Width,Height:Height,AppliedMutations:[]CandidateMutation{mutation}}, nil
}

// drawVisibleRuleRoleBindingR1 separates rule operands by semantic role so a
// receiver does not need to infer which cell is the source condition, which is
// the target, or which target state is required before applying the SET.
func drawVisibleRuleRoleBindingR1(img *image.Gray, decoded TemporalCarrierDecoded) {
	labels := temporalVisibleCellLabels(decoded, true)
	p := decoded.Program
	if len(p.Automaton.Cells) > 4 {
		drawText(img, 24, 300, 1, "RULE ROLE BINDING: DENSE T2 FALLBACK", 0)
		emphasizeTemporalCheckpoints(img, decoded)
		return
	}

	fill(img, 18, 202, 604, 112, 0xff)
	box(img, 18, 202, 604, 112, 0)
	drawText(img, 24, 208, 1, "T2 RULE MICROGRAMMAR | EXPLICIT OPERAND ROLES", 0)
	drawText(img, 24, 220, 1, "EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT", 0)

	y := 234
	visible := 0
	for _, r := range p.Automaton.Rules {
		if visible >= 4 { break }
		req := "TRUE"
		if len(r.Requires) > 0 {
			parts := make([]string, 0, len(r.Requires))
			for _, pred := range r.Requires {
				parts = append(parts, shortLabel(labels[pred.CellID], 16)+"="+shortLabel(pred.State, 10))
			}
			req = strings.Join(parts, "&")
		}
		from := r.FromState
		if strings.TrimSpace(from) == "" { from = "*" }
		drawText(img, 28, y, 1, shortLabel(fmt.Sprintf("%s WHEN %s", strings.ToUpper(r.ID), req), 94), 0)
		y += 10
		drawText(img, 36, y, 1, shortLabel(fmt.Sprintf("TARGET %s | REQUIRE %s | SET %s", shortLabel(labels[r.TargetCell], 16), shortLabel(from, 10), shortLabel(r.ToState, 10)), 94), 0)
		y += 12
		visible++
	}
	if len(p.Automaton.Rules) > visible {
		drawText(img, 28, y, 1, fmt.Sprintf("+%d MORE RULES IN EXACT PROGRAM", len(p.Automaton.Rules)-visible), 0)
	}
}
