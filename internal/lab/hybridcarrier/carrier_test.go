package hybridcarrier

import (
	"encoding/json"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/receiver"
)

func receiverSpec() receiver.Spec {
	return receiver.Spec{
		ContractID: receiver.ContractID,
		Boot: receiver.BootSpec{
			Version:          "r0",
			EntrySymbol:      "BOOT",
			StartState:       "s0",
			RosettaRef:       "ROSETTA",
			ProgramRef:       "PROGRAM",
			IndexRef:         "INDEX",
			MemoryRef:        "MEMORY",
			VerificationRef:  "VERIFICATION",
			RequiredSections: []string{"BOOT", "ROSETTA", "PROGRAM", "INDEX", "MEMORY", "VERIFICATION"},
			WorkingWindow:    4000,
		},
		Rosetta: []receiver.SymbolBinding{
			{Physical: "triangle", Semantic: "OPEN"},
			{Physical: "circle", Semantic: "VALUE"},
			{Physical: "square", Semantic: "HALT"},
		},
		Program: []receiver.MicroRule{
			{ID: "m0", State: "s0", Token: "OPEN", NextState: "s1"},
			{ID: "m1", State: "s1", Token: "VALUE", Emit: "READ", NextState: "done"},
			{ID: "m2", State: "done", Token: "HALT", NextState: "done", Halt: true},
		},
	}
}

func TestEnvelopeRoundTripAndMemoryVerification(t *testing.T) {
	envelope, err := BuildEnvelope(
		receiverSpec(),
		[]IndexEntry{{Key: "K7F91", Address: "mem:1"}},
		[]MemoryEntry{
			{Address: "mem:1", Value: "AMBER-10593", Links: map[string][]string{"depends": []string{"mem:2"}}},
			{Address: "mem:2", Value: "DELTA-77102"},
		},
	)
	if err != nil { t.Fatal(err) }
	payload, err := EncodePayload(envelope)
	if err != nil { t.Fatal(err) }
	decoded, err := DecodePayload(payload)
	if err != nil { t.Fatal(err) }
	if decoded.Memory[0].Address != "mem:1" || decoded.Memory[0].Value != "AMBER-10593" {
		t.Fatalf("unexpected decoded memory: %+v", decoded.Memory)
	}

	decoded.Memory[0].Value = "MUTATED"
	mutated, err := json.Marshal(decoded)
	if err != nil { t.Fatal(err) }
	mutated = append([]byte(MagicR0), mutated...)
	if _, err := DecodePayload(mutated); err == nil {
		t.Fatal("expected memory verification to detect mutation")
	}
}

func TestGlyphTransportRoundTripAndPNGTarget(t *testing.T) {
	envelope, err := BuildEnvelope(
		receiverSpec(),
		[]IndexEntry{{Key: "K7F91", Address: "mem:1"}},
		[]MemoryEntry{{Address: "mem:1", Value: "AMBER-10593"}},
	)
	if err != nil { t.Fatal(err) }
	payload, err := EncodePayload(envelope)
	if err != nil { t.Fatal(err) }
	glyphs := glyphcalc.GlyphsFromBytes(payload)
	recovered, err := glyphcalc.BytesFromGlyphs(glyphs)
	if err != nil { t.Fatal(err) }
	if string(recovered) != string(payload) {
		t.Fatal("glyph transport changed carrier payload")
	}

	rendered, err := Render(envelope, 24, 40, 8, DefaultMaxPNGBytes)
	if err != nil { t.Fatal(err) }
	if rendered.PNGBytes == 0 || rendered.PNGBytes > DefaultMaxPNGBytes {
		t.Fatalf("unexpected PNG size %d", rendered.PNGBytes)
	}
	if rendered.PayloadSHA256 == "" || rendered.PNGSHA256 == "" {
		t.Fatal("missing carrier identities")
	}
}
