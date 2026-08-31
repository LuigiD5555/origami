package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	stddraw "image/draw"
	"image/png"
)

const ExecuteDontSummarizeToStableR1 = "EXECUTE_DONT_SUMMARIZE_TO_STABLE_R1"

const (
	ExecutionComplianceModeTextR1  = "MODE: EXECUTE PROGRAM | DO NOT SUMMARIZE OR DESCRIBE"
	ExecutionComplianceStep1TextR1 = "1 READ INITIAL CELL STATES"
	ExecutionComplianceStep2TextR1 = "2 TEST ALL RULES"
	ExecutionComplianceStep3TextR1 = "3 USE SAME PRE-STEP"
	ExecutionComplianceStep4TextR1 = "4 APPLY FIRED SETS"
	ExecutionComplianceStep5TextR1 = "5 REPEAT UNTIL UNCHANGED"
	ExecutionComplianceStep6TextR1 = "6 REPORT FINAL A/B/C"
	ExecutionComplianceStopTextR1  = "STOP ONLY WHEN UNCHANGED"
)

// BuildExecutionComplianceCandidate is the deterministic R6 materializer. It
// modifies only the visible execution-policy surface. The exact TemporalProgram
// and the R5 rule-role grammar remain untouched and are verified after render.
func BuildExecutionComplianceCandidate(parentPNG []byte, id string) ([]byte, CandidateBuildReport, error) {
	parentSum := sha256.Sum256(parentPNG)
	parentSHA := hex.EncodeToString(parentSum[:])
	parentDecoded, err := DecodeTemporalCarrierPNG(parentPNG)
	if err != nil { return nil, CandidateBuildReport{}, fmt.Errorf("parent is not a valid temporal carrier: %w", err) }
	img, err := png.Decode(bytes.NewReader(parentPNG))
	if err != nil { return nil, CandidateBuildReport{}, err }
	gray := image.NewGray(image.Rect(0, 0, Width, Height))
	stddraw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, stddraw.Src)
	drawExecutionComplianceR1(gray)

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
	mutation := CandidateMutation{Kind:"PROMPT",Target:"EXECUTION_POLICY_COMPLIANCE",Value:ExecuteDontSummarizeToStableR1,Experimental:true}
	return candidate, CandidateBuildReport{
		Schema:CandidateBuildReportSchema, CandidateID:id, ParentProfile:parentDecoded.Profile,
		ParentSHA256:parentSHA, CandidateSHA256:hex.EncodeToString(candidateSum[:]),
		ParentProgramSHA256:parentDecoded.ProgramSHA256, CandidateProgramSHA256:candidateDecoded.ProgramSHA256,
		ExactProgramPreserved:true, PNGBytes:len(candidate), Width:Width, Height:Height,
		AppliedMutations:[]CandidateMutation{mutation},
	}, nil
}

func drawExecutionComplianceR1(img *image.Gray) {
	// Replace only the previous compact execution-policy row with an imperative
	// mode declaration. This stays inside the BOOT guidance surface.
	fill(img, 20, 88, 596, 14, 0xff)
	drawText(img, 24, 92, 1, ExecutionComplianceModeTextR1, 0)

	// The exact payload starts at x=224. The left margin is intentionally used as
	// an execution checklist so the payload, rule grammar and timeline stay frozen.
	fill(img, 18, 420, 190, 154, 0xff)
	box(img, 18, 420, 190, 154, 0)
	drawText(img, 26, 428, 1, "EXECUTION REQUIRED", 0)
	drawText(img, 26, 444, 1, ExecutionComplianceStep1TextR1, 0)
	drawText(img, 26, 458, 1, ExecutionComplianceStep2TextR1, 0)
	drawText(img, 26, 472, 1, ExecutionComplianceStep3TextR1, 0)
	drawText(img, 26, 486, 1, ExecutionComplianceStep4TextR1, 0)
	drawText(img, 26, 500, 1, ExecutionComplianceStep5TextR1, 0)
	drawText(img, 26, 514, 1, ExecutionComplianceStep6TextR1, 0)
	drawText(img, 26, 536, 1, ExecutionComplianceStopTextR1, 0)
}
