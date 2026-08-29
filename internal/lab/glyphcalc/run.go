package glyphcalc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Metrics struct {
	Schema                 string  `json:"schema"`
	Pattern                string  `json:"pattern"`
	Seed                   uint64  `json:"seed"`
	SourceBytes            int     `json:"source_bytes"`
	ProgramBytesCandidate  int     `json:"program_bytes_candidate"`
	SelectedMode           string  `json:"selected_mode"`
	TransportBytes         int     `json:"transport_bytes"`
	GlyphCount             int     `json:"glyph_count"`
	InstructionsExecuted   int     `json:"instructions_executed"`
	GlyphsPerSourceByte    float64 `json:"glyphs_per_source_byte"`
	SourceBytesPerGlyph    float64 `json:"source_bytes_per_glyph"`
	PerceptionReductionPct float64 `json:"perception_reduction_pct"`
	ExactReconstruction    bool    `json:"exact_reconstruction"`
	SourceSHA256           string  `json:"source_sha256"`
	TransportSHA256        string  `json:"transport_sha256"`
	ReconstructionSHA256   string  `json:"reconstruction_sha256"`
	CarrierPNGBytes        int     `json:"carrier_png_bytes"`
	CarrierPNGSHA256       string  `json:"carrier_png_sha256"`
	CarrierPixelSHA256     string  `json:"carrier_pixel_sha256"`
	CompilerProgramStats   Stats   `json:"compiler_program_stats"`
	ExecutedProgramStats   Stats   `json:"executed_program_stats"`
}

type GroundTruth struct {
	Schema           string  `json:"schema"`
	SemanticMeaning  string  `json:"semantic_meaning"`
	PhysicalCodebook string  `json:"physical_codebook"`
	Formula          string  `json:"formula"`
	TransportMode    string  `json:"transport_mode"`
	Glyphs           []Glyph `json:"glyphs"`
}

func Run(outDir, pattern string, length int, seed uint64) (Metrics, error) {
	src, err := Generate(pattern, length, seed)
	if err != nil {
		return Metrics{}, err
	}
	tr, err := CompileBest(src)
	if err != nil {
		return Metrics{}, err
	}
	recon, execStats, err := ExecuteTransport(tr.Stream, max(length*2, 1024), max(len(tr.Stream)*2, 128))
	if err != nil {
		return Metrics{}, err
	}
	exact := string(src) == string(recon)
	if !exact {
		return Metrics{}, fmt.Errorf("internal exact reconstruction failed")
	}
	gs := GlyphsFromBytes(tr.Stream)
	pngb, rmeta, err := RenderGlyphGrid(gs, 16, 52, 12)
	if err != nil {
		return Metrics{}, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return Metrics{}, err
	}
	files := map[string][]byte{}
	files["source.bin"] = src
	files["transport.bin"] = tr.Stream
	files["program_candidate.bin"] = tr.Program
	files["reconstruction.bin"] = recon
	files["carrier.png"] = pngb
	gt := GroundTruth{Schema: "ohf.glyph-calculus-ground-truth.v0", SemanticMeaning: "NONE", PhysicalCodebook: "R3.10 F01xF02xF03xF04", Formula: "byte = shape + 4*holes + 16*direction + 64*frames", TransportMode: tr.Mode, Glyphs: gs}
	gtb, _ := json.MarshalIndent(gt, "", "  ")
	gtb = append(gtb, '\n')
	files["glyphs.json"] = gtb
	for n, b := range files {
		if err := os.WriteFile(filepath.Join(outDir, n), b, 0o644); err != nil {
			return Metrics{}, err
		}
	}
	m := Metrics{Schema: "ohf.glyph-calculus-metrics.v0", Pattern: pattern, Seed: seed, SourceBytes: len(src), ProgramBytesCandidate: len(tr.Program) + 1, SelectedMode: tr.Mode, TransportBytes: len(tr.Stream), GlyphCount: len(gs), InstructionsExecuted: execStats.Instructions, GlyphsPerSourceByte: float64(len(gs)) / float64(len(src)), SourceBytesPerGlyph: float64(len(src)) / float64(len(gs)), PerceptionReductionPct: (1 - float64(len(gs))/float64(len(src))) * 100, ExactReconstruction: exact, SourceSHA256: hash(src), TransportSHA256: hash(tr.Stream), ReconstructionSHA256: hash(recon), CarrierPNGBytes: rmeta.PNGBytes, CarrierPNGSHA256: rmeta.PNGSHA256, CarrierPixelSHA256: rmeta.PixelSHA256, CompilerProgramStats: tr.ProgramStat, ExecutedProgramStats: execStats}
	mb, _ := json.MarshalIndent(m, "", "  ")
	mb = append(mb, '\n')
	if err := os.WriteFile(filepath.Join(outDir, "metrics.json"), mb, 0o644); err != nil {
		return Metrics{}, err
	}
	rmb, _ := json.MarshalIndent(rmeta, "", "  ")
	rmb = append(rmb, '\n')
	if err := os.WriteFile(filepath.Join(outDir, "render.json"), rmb, 0o644); err != nil {
		return Metrics{}, err
	}
	if err := writeNativeFiles(outDir, tr.Stream); err != nil {
		return Metrics{}, err
	}
	if err := writeChecksums(outDir); err != nil {
		return Metrics{}, err
	}
	return m, nil
}

func Verify(outDir string) error {
	must := []string{"source.bin", "transport.bin", "reconstruction.bin", "carrier.png", "glyphs.json", "metrics.json", "render.json", "CHECKSUMS.sha256"}
	for _, n := range must {
		if _, err := os.Stat(filepath.Join(outDir, n)); err != nil {
			return fmt.Errorf("missing %s: %w", n, err)
		}
	}
	if err := verifyChecksums(outDir); err != nil {
		return err
	}
	src, _ := os.ReadFile(filepath.Join(outDir, "source.bin"))
	stream, _ := os.ReadFile(filepath.Join(outDir, "transport.bin"))
	recon, _, err := ExecuteTransport(stream, max(len(src)*2, 1024), max(len(stream)*2, 128))
	if err != nil {
		return err
	}
	if string(src) != string(recon) {
		return fmt.Errorf("reconstruction mismatch")
	}
	return nil
}

func hash(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func writeNativeFiles(dir string, stream []byte) error {
	native := filepath.Join(dir, "native")
	if err := os.MkdirAll(native, 0o755); err != nil {
		return err
	}
	prompt := `# OHF Glyph Calculus — Native Perception Probe

You are NOT asked to execute, decompress, infer meaning, or reason about the payload.
Your only task is local visual perception.

Read glyphs row-major (left-to-right, top-to-bottom).
Each glyph has exactly four physical states:
- outer shape: circle=0, square=1, diamond=2, triangle=3
- holes: 0,1,2,3
- direction beacon: top=0, right=1, bottom=2, left=3
- frame count: 0,1,2,3

For each glyph compute only:
byte = shape + 4*holes + 16*direction + 64*frame_count

Return ONLY JSON with this schema:
{"status":"KNOWN|AMBIGUOUS|UNKNOWN|INVALID","bytes_hex":"lowercase hexadecimal with no spaces","glyph_count":integer}

Do not execute the byte stream. Do not assign semantic meaning to shapes. If any glyph is not uniquely readable, use AMBIGUOUS or UNKNOWN rather than guessing.
`
	query := `Decode every visible glyph into bytes using only the local physical rule. Return the requested JSON and do nothing else.`
	expected := fmt.Sprintf("{\n  \"status\": \"KNOWN\",\n  \"bytes_hex\": \"%s\",\n  \"glyph_count\": %d\n}\n", hex.EncodeToString(stream), len(stream))
	for n, s := range map[string]string{"MASTER_PROMPT.md": prompt, "query.txt": query, "expected_observation.json": expected} {
		if err := os.WriteFile(filepath.Join(native, n), []byte(s), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeChecksums(dir string) error {
	var names []string
	ents, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range ents {
		if !e.IsDir() && e.Name() != "CHECKSUMS.sha256" {
			names = append(names, e.Name())
		}
	}
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[i] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	var out []byte
	for _, n := range names {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			return err
		}
		out = append(out, []byte(fmt.Sprintf("%s  %s\n", hash(b), n))...)
	}
	return os.WriteFile(filepath.Join(dir, "CHECKSUMS.sha256"), out, 0o644)
}

func verifyChecksums(dir string) error {
	b, err := os.ReadFile(filepath.Join(dir, "CHECKSUMS.sha256"))
	if err != nil {
		return err
	}
	var wantHash, wantName string
	for _, line := range splitLines(string(b)) {
		if line == "" {
			continue
		}
		if _, err := fmt.Sscanf(line, "%s  %s", &wantHash, &wantName); err != nil {
			return fmt.Errorf("bad checksum line %q", line)
		}
		fb, err := os.ReadFile(filepath.Join(dir, wantName))
		if err != nil {
			return err
		}
		if hash(fb) != wantHash {
			return fmt.Errorf("checksum mismatch %s", wantName)
		}
	}
	return nil
}
func splitLines(s string) []string {
	var out []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
