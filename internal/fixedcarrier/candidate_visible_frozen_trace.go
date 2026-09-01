package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"

	"github.com/LuigiD5555/origami/internal/temporal"
)

// BuildVisibleFrozenTraceCandidate extends a temporal carrier by rendering
// frozen checkpoint states below the timeline. Each checkpoint is rendered as
// a compact row of per-cell state glyphs (O/*/X for IDLE/ACTIVE/DONE), so a
// small VLM can read already-solved computation instead of simulating the
// program. The exact-plane payload (SHA-256+zlib+CRC) is untouched and its
// program_sha256 must remain identical to the parent.
func BuildVisibleFrozenTraceCandidate(parentPNG []byte, id string) ([]byte, CandidateBuildReport, error) {
	parentSum := sha256.Sum256(parentPNG)
	parentSHA := hex.EncodeToString(parentSum[:])
	parentDecoded, err := DecodeTemporalCarrierPNG(parentPNG)
	if err != nil {
		return nil, CandidateBuildReport{}, fmt.Errorf("parent is not a valid temporal carrier: %w", err)
	}

	// Run the program to obtain the trace with frozen checkpoint states.
	trace, err := temporal.Run(parentDecoded.Program)
	if err != nil {
		return nil, CandidateBuildReport{}, fmt.Errorf("run temporal program for frozen trace: %w", err)
	}

	// Re-render the carrier with the trace, so drawTemporalSemanticPlane
	// can access checkpoint states and render glifos.
	img := image.NewGray(image.Rect(0, 0, Width, Height))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	drawFrame(img)
	drawTemporalBoot(img)

	// Pass the trace to enable frozen checkpoint state rendering.
	drawTemporalSemanticPlane(img, parentDecoded.Program, &trace)

	// Reconstruct the exact-plane from the original program.
	canonical, err := json.Marshal(parentDecoded.Program)
	if err != nil {
		return nil, CandidateBuildReport{}, err
	}
	digest := sha256.Sum256(canonical)
	compressed, err := zlibBest(canonical)
	if err != nil {
		return nil, CandidateBuildReport{}, err
	}
	const headerBytes = 46
	const crcBytes = 4
	capacity := RecordBytes - headerBytes - crcBytes
	if len(compressed) > capacity {
		return nil, CandidateBuildReport{}, fmt.Errorf("program compressed size exceeds capacity: %d > %d", len(compressed), capacity)
	}
	rec := make([]byte, RecordBytes)
	copy(rec[:8], temporalMagic[:])
	binary.BigEndian.PutUint32(rec[8:12], uint32(len(canonical)))
	binary.BigEndian.PutUint16(rec[12:14], uint16(len(compressed)))
	copy(rec[14:46], digest[:])
	copy(rec[46:46+len(compressed)], compressed)
	binary.BigEndian.PutUint32(rec[508:512], crc32.ChecksumIEEE(rec[:508]))

	drawTemporalExactPlane(img, rec)

	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&raw, img); err != nil {
		return nil, CandidateBuildReport{}, err
	}
	if raw.Len() > FixedPNGBytes-12 {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate raw PNG exceeds frozen envelope: %d > %d", raw.Len(), FixedPNGBytes-12)
	}
	candidate, err := padPNG(raw.Bytes(), FixedPNGBytes)
	if err != nil {
		return nil, CandidateBuildReport{}, err
	}

	// Verify exact-plane: program_sha256 must be identical.
	candidateDecoded, err := DecodeTemporalCarrierPNG(candidate)
	if err != nil {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate exact-plane verification failed: %w", err)
	}
	if candidateDecoded.ProgramSHA256 != parentDecoded.ProgramSHA256 {
		return nil, CandidateBuildReport{}, fmt.Errorf("candidate mutated exact TemporalProgram")
	}

	candidateSum := sha256.Sum256(candidate)
	mutation := CandidateMutation{
		Kind:         "PROMPT",
		Target:       "TEMPORAL_STRUCTURE",
		Value:        "VISIBLE_FROZEN_TRACE_R1",
		Experimental: true,
	}
	return candidate, CandidateBuildReport{
		Schema:                 CandidateBuildReportSchema,
		CandidateID:            id,
		ParentProfile:          parentDecoded.Profile,
		ParentSHA256:           parentSHA,
		CandidateSHA256:        hex.EncodeToString(candidateSum[:]),
		ParentProgramSHA256:    parentDecoded.ProgramSHA256,
		CandidateProgramSHA256: candidateDecoded.ProgramSHA256,
		ExactProgramPreserved:  true,
		PNGBytes:               len(candidate),
		Width:                  Width,
		Height:                 Height,
		AppliedMutations:       []CandidateMutation{mutation},
	}, nil
}
