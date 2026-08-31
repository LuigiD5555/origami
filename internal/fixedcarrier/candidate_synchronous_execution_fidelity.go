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
	SyncFidelityHeaderTextR1 = "SYNC: FREEZE > SELECT > APPLY"
	SyncFidelityStep1TextR1  = "1 FREEZE PRE-STEP"
	SyncFidelityStep2TextR1  = "2 FIRE IFF WHEN+TARGET=REQUIRE"
	SyncFidelityStep3TextR1  = "3 SELECT ALL; APPLY TOGETHER"
	SyncFidelityNoOrderTextR1 = "4 NO ORDER; NO CASCADE"
	SyncFidelityStopTextR1 = "5 NEXT=UPDATED; STOP IF ZERO FIRE"
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
	fill(img, 18, 420, 190, 92, 0xff)
	box(img, 18, 420, 190, 92, 0)
	drawText(img, 24, 426, 1, SyncFidelityHeaderTextR1, 0)
	drawText(img, 24, 442, 1, SyncFidelityStep1TextR1, 0)
	drawText(img, 24, 456, 1, SyncFidelityStep2TextR1, 0)
	drawText(img, 24, 470, 1, SyncFidelityStep3TextR1, 0)
	drawText(img, 24, 484, 1, SyncFidelityNoOrderTextR1, 0)
	drawText(img, 24, 498, 1, SyncFidelityStopTextR1, 0)
}
