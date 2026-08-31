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

const FreezeSelectApplyTogetherR1 = "FREEZE_SELECT_APPLY_TOGETHER_R1"

const (
	SyncFidelityHeaderTextR1 = "SYNC STEP = FREEZE > SELECT > APPLY"
	SyncFidelityStep1TextR1  = "1 FREEZE PRE-STEP SNAPSHOT"
	SyncFidelityStep2TextR1  = "2 RULE FIREABLE ONLY IF"
	SyncFidelityStep2BTextR1 = "  WHEN TRUE + TARGET=REQUIRE"
	SyncFidelityStep3TextR1  = "3 SELECT ALL FIREABLE RULES"
	SyncFidelityStep4TextR1  = "4 APPLY ALL SETS TOGETHER"
	SyncFidelityNoOrderTextR1 = "NO RULE ORDER"
	SyncFidelityNoCascadeTextR1 = "NO CASCADE INSIDE STEP"
	SyncFidelityNextTextR1 = "NEXT USES UPDATED SNAPSHOT"
	SyncFidelityStopTextR1 = "STOP WHEN ZERO RULES FIRE"
)

// BuildSynchronousExecutionFidelityCandidate is the deterministic R7 materializer.
// It changes only the visible synchronous-execution guidance inherited from R6.
// Rule grammar, states, program bytes and exact payload remain frozen.
func BuildSynchronousExecutionFidelityCandidate(parentPNG []byte, id string) ([]byte, CandidateBuildReport, error) {
	parentSum := sha256.Sum256(parentPNG)
	parentSHA := hex.EncodeToString(parentSum[:])
	parentDecoded, err := DecodeTemporalCarrierPNG(parentPNG)
	if err != nil { return nil, CandidateBuildReport{}, fmt.Errorf("parent is not a valid temporal carrier: %w", err) }
	img, err := png.Decode(bytes.NewReader(parentPNG))
	if err != nil { return nil, CandidateBuildReport{}, err }
	gray := image.NewGray(image.Rect(0, 0, Width, Height))
	stddraw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, stddraw.Src)
	drawSynchronousExecutionFidelityR1(gray)

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
	mutation := CandidateMutation{Kind:"PROMPT",Target:"SYNCHRONOUS_EXECUTION_FIDELITY",Value:FreezeSelectApplyTogetherR1,Experimental:true}
	return candidate, CandidateBuildReport{
		Schema:CandidateBuildReportSchema, CandidateID:id, ParentProfile:parentDecoded.Profile,
		ParentSHA256:parentSHA, CandidateSHA256:hex.EncodeToString(candidateSum[:]),
		ParentProgramSHA256:parentDecoded.ProgramSHA256, CandidateProgramSHA256:candidateDecoded.ProgramSHA256,
		ExactProgramPreserved:true, PNGBytes:len(candidate), Width:Width, Height:Height,
		AppliedMutations:[]CandidateMutation{mutation},
	}, nil
}

func drawSynchronousExecutionFidelityR1(img *image.Gray) {
	// R7 intentionally reuses only R6's left execution guidance panel. The
	// imperative BOOT row, T2 microgrammar, timeline and exact payload stay frozen.
	fill(img, 18, 420, 190, 154, 0xff)
	box(img, 18, 420, 190, 154, 0)
	drawText(img, 24, 426, 1, SyncFidelityHeaderTextR1, 0)
	drawText(img, 24, 440, 1, SyncFidelityStep1TextR1, 0)
	drawText(img, 24, 453, 1, SyncFidelityStep2TextR1, 0)
	drawText(img, 24, 466, 1, SyncFidelityStep2BTextR1, 0)
	drawText(img, 24, 479, 1, SyncFidelityStep3TextR1, 0)
	drawText(img, 24, 492, 1, SyncFidelityStep4TextR1, 0)
	drawText(img, 24, 508, 1, SyncFidelityNoOrderTextR1, 0)
	drawText(img, 24, 521, 1, SyncFidelityNoCascadeTextR1, 0)
	drawText(img, 24, 538, 1, SyncFidelityNextTextR1, 0)
	drawText(img, 24, 555, 1, SyncFidelityStopTextR1, 0)
}
