package hybridcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/receiver"
)

const (
	SchemaR0            = "origami.hybrid-carrier.r0"
	MagicR0             = "ORIGAMI-HYBRID-R0\n"
	DefaultMaxPNGBytes  = 512000
)

type IndexEntry struct {
	Key     string `json:"key"`
	Address string `json:"address"`
}

type MemoryEntry struct {
	Address string            `json:"address"`
	Value   string            `json:"value"`
	Links   map[string][]string `json:"links,omitempty"`
}

type Verification struct {
	MemorySHA256 string `json:"memory_sha256"`
}

type Envelope struct {
	Schema       string            `json:"schema"`
	Receiver     receiver.Spec     `json:"receiver"`
	Index        []IndexEntry      `json:"index"`
	Memory       []MemoryEntry     `json:"memory"`
	Verification Verification      `json:"verification"`
}

type Rendered struct {
	PNG              []byte                `json:"-"`
	PNGBytes         int                   `json:"png_bytes"`
	PNGSHA256        string                `json:"png_sha256"`
	PayloadBytes     int                   `json:"payload_bytes"`
	PayloadSHA256    string                `json:"payload_sha256"`
	GlyphCount       int                   `json:"glyph_count"`
	RenderMeta       glyphcalc.RenderMeta  `json:"render_meta"`
}

func BuildEnvelope(spec receiver.Spec, index []IndexEntry, memory []MemoryEntry) (Envelope, error) {
	if err := receiver.Validate(spec); err != nil {
		return Envelope{}, fmt.Errorf("receiver spec: %w", err)
	}
	idx := append([]IndexEntry(nil), index...)
	mem := append([]MemoryEntry(nil), memory...)
	sort.Slice(idx, func(i, j int) bool {
		if idx[i].Key == idx[j].Key { return idx[i].Address < idx[j].Address }
		return idx[i].Key < idx[j].Key
	})
	sort.Slice(mem, func(i, j int) bool { return mem[i].Address < mem[j].Address })
	if len(mem) == 0 {
		return Envelope{}, fmt.Errorf("memory cannot be empty")
	}
	seen := map[string]bool{}
	for _, entry := range mem {
		if entry.Address == "" { return Envelope{}, fmt.Errorf("memory address cannot be empty") }
		if seen[entry.Address] { return Envelope{}, fmt.Errorf("duplicate memory address %q", entry.Address) }
		seen[entry.Address] = true
	}
	for _, entry := range idx {
		if entry.Key == "" || entry.Address == "" { return Envelope{}, fmt.Errorf("index entries require key and address") }
		if !seen[entry.Address] { return Envelope{}, fmt.Errorf("index key %q points to missing address %q", entry.Key, entry.Address) }
	}
	memoryBytes, err := json.Marshal(mem)
	if err != nil { return Envelope{}, err }
	return Envelope{
		Schema:   SchemaR0,
		Receiver: spec,
		Index:    idx,
		Memory:   mem,
		Verification: Verification{MemorySHA256: hash(memoryBytes)},
	}, nil
}

func EncodePayload(envelope Envelope) ([]byte, error) {
	if envelope.Schema != SchemaR0 { return nil, fmt.Errorf("carrier schema must be %q", SchemaR0) }
	if err := receiver.Validate(envelope.Receiver); err != nil { return nil, err }
	body, err := json.Marshal(envelope)
	if err != nil { return nil, err }
	payload := make([]byte, 0, len(MagicR0)+len(body))
	payload = append(payload, []byte(MagicR0)...)
	payload = append(payload, body...)
	return payload, nil
}

func DecodePayload(payload []byte) (Envelope, error) {
	if !bytes.HasPrefix(payload, []byte(MagicR0)) {
		return Envelope{}, fmt.Errorf("missing Hybrid carrier magic")
	}
	var envelope Envelope
	if err := json.Unmarshal(payload[len(MagicR0):], &envelope); err != nil { return Envelope{}, err }
	if envelope.Schema != SchemaR0 { return Envelope{}, fmt.Errorf("unexpected carrier schema %q", envelope.Schema) }
	if err := receiver.Validate(envelope.Receiver); err != nil { return Envelope{}, err }
	memoryBytes, err := json.Marshal(envelope.Memory)
	if err != nil { return Envelope{}, err }
	if got := hash(memoryBytes); got != envelope.Verification.MemorySHA256 {
		return Envelope{}, fmt.Errorf("memory verification mismatch: expected=%s actual=%s", envelope.Verification.MemorySHA256, got)
	}
	return envelope, nil
}

// Render projects the complete logical carrier payload into the existing
// Glyph Calculus physical transport. This proves deterministic integration and
// exact round-trip at the transport layer; it does NOT by itself prove that a
// VLM can perceptually recover BOOT/ROSETTA from the image.
func Render(envelope Envelope, columns, cell, margin, maxPNGBytes int) (Rendered, error) {
	payload, err := EncodePayload(envelope)
	if err != nil { return Rendered{}, err }
	glyphs := glyphcalc.GlyphsFromBytes(payload)
	pngBytes, meta, err := glyphcalc.RenderGlyphGrid(glyphs, columns, cell, margin)
	if err != nil { return Rendered{}, err }
	if maxPNGBytes <= 0 { maxPNGBytes = DefaultMaxPNGBytes }
	if len(pngBytes) > maxPNGBytes {
		return Rendered{}, fmt.Errorf("carrier exceeds PNG target: bytes=%d max=%d", len(pngBytes), maxPNGBytes)
	}
	return Rendered{
		PNG:           pngBytes,
		PNGBytes:      len(pngBytes),
		PNGSHA256:     hash(pngBytes),
		PayloadBytes:  len(payload),
		PayloadSHA256: hash(payload),
		GlyphCount:    len(glyphs),
		RenderMeta:    meta,
	}, nil
}

func hash(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
