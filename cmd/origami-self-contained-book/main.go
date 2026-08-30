package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LuigiD5555/origami/internal/document"
	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/integration"
	"github.com/LuigiD5555/origami/internal/lab/selfcontainedbook"
)

type report struct {
	Schema              string                   `json:"schema"`
	Profile             string                   `json:"profile"`
	Status              string                   `json:"status"`
	SourcePath          string                   `json:"source_path,omitempty"`
	SourcePlane         string                   `json:"source_plane"`
	Header              selfcontainedbook.Header `json:"header"`
	Control             fixedcarrier.Decoded     `json:"control"`
	PNGBytes            int                      `json:"png_bytes"`
	TargetPNGBytes      int                      `json:"target_png_bytes"`
	OverTargetBytes     int                      `json:"over_target_bytes"`
	DataWidthBits       int                      `json:"data_width_bits"`
	DataHeightBits      int                      `json:"data_height_bits"`
	RoundtripExact      bool                     `json:"roundtrip_exact"`
	SemanticFoldClaim   string                   `json:"semantic_fold_claim"`
	ExactResidualClaim  string                   `json:"exact_residual_claim"`
	MasterPromptPath    string                   `json:"master_prompt_path"`
	MasterPromptSHA256  string                   `json:"master_prompt_sha256"`
	ModelPacketPath     string                   `json:"model_packet_path"`
}

type modelPacket struct {
	Schema                 string   `json:"schema"`
	Profile                string   `json:"profile"`
	Carrier                string   `json:"carrier"`
	CarrierSHA256          string   `json:"carrier_sha256"`
	MasterPrompt           string   `json:"master_prompt"`
	MasterPromptSHA256     string   `json:"master_prompt_sha256"`
	SourceSHA256           string   `json:"source_sha256"`
	WorkingWindowTokenEq   int      `json:"working_window_token_eq"`
	RequiredInputs         []string `json:"required_inputs"`
	ReadOrder              []string `json:"read_order"`
}

func main() {
	mode := flag.String("mode", "build", "build|decode")
	in := flag.String("in", "", "PDF input for build or carrier PNG for decode")
	out := flag.String("out", "origami-self-contained-book-r0.png", "carrier PNG for build or recovered UTF-8 text for decode")
	prompt := flag.String("prompt", "generated/MASTER_PROMPT.md", "Origami Master Prompt used/bound by the carrier")
	dataWidth := flag.Int("data-width", selfcontainedbook.DefaultDataWidth, "exact residual matrix width in bits")
	flag.Parse()
	if *in == "" { die(fmt.Errorf("-in is required")) }
	switch *mode {
	case "build": die(build(*in, *out, *prompt, *dataWidth))
	case "decode": die(decode(*in, *out))
	default: die(fmt.Errorf("unsupported mode %q", *mode))
	}
}

func build(path, out, promptPath string, dataWidth int) error {
	promptBytes, err := os.ReadFile(promptPath); if err != nil { return fmt.Errorf("read Master Prompt: %w", err) }
	promptSHA := hashBytes(promptBytes)
	text, err := extractPDFText(path); if err != nil { return err }
	doc, err := document.FromText(filepath.Base(path), "application/pdf", text); if err != nil { return err }
	bundle, err := document.BuildBundle(doc); if err != nil { return err }
	system, err := integration.NewSystem(doc.SourceID, bundle.Nodes); if err != nil { return err }
	blockCount := 0; for _, page := range doc.Pages { blockCount += len(page.Blocks) }
	signature, err := json.Marshal(system.Signature()); if err != nil { return err }
	compressed, err := bzip2Compress(text); if err != nil { return err }
	payload, header, err := selfcontainedbook.Pack(text, compressed, uint32(len(doc.Pages)), uint32(blockCount), "bz2-9", promptSHA); if err != nil { return err }
	meta := fixedcarrier.Metadata{CarrierID: doc.SourceHash, StoreRoot: integration.CanonicalStoreRoot(system.Nodes()), SourceSHA256: doc.SourceHash, PageCount: uint32(len(doc.Pages)), BlockCount: uint32(blockCount), DocumentCount: 1, ObjectCount: uint32(len(bundle.Nodes)), GraphSignature: signature}
	controlPNG, control, err := fixedcarrier.Render(meta); if err != nil { return err }
	rendered, err := selfcontainedbook.Render(controlPNG, payload, dataWidth); if err != nil { return err }
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil { return err }
	if err := os.WriteFile(out, rendered.PNG, 0o644); err != nil { return err }

	// Portable model-facing package. The Markdown file is the actual prompt;
	// model_packet.json only binds the exact prompt/carrier pair and read order.
	packageDir := filepath.Dir(out)
	masterPromptOut := filepath.Join(packageDir, "MASTER_PROMPT.md")
	modelPacketOut := filepath.Join(packageDir, "model_packet.json")
	if err := os.WriteFile(masterPromptOut, promptBytes, 0o644); err != nil { return err }
	packet := modelPacket{
		Schema: "origami.self-contained-book.r0-lab.model-packet",
		Profile: selfcontainedbook.ProfileID,
		Carrier: filepath.Base(out), CarrierSHA256: hashBytes(rendered.PNG),
		MasterPrompt: filepath.Base(masterPromptOut), MasterPromptSHA256: promptSHA,
		SourceSHA256: header.SourceSHA256, WorkingWindowTokenEq: 4000,
		RequiredInputs: []string{"MASTER_PROMPT.md", filepath.Base(out), "user question"},
		ReadOrder: []string{"MASTER_PROMPT", "T0_BOOT", "T1_ROSETTA", "T2_SUPERINDEX", "T3_MEMORY", "VERIFY"},
	}
	packetBytes, err := json.MarshalIndent(packet, "", "  "); if err != nil { return err }
	packetBytes = append(packetBytes, '\n'); if err := os.WriteFile(modelPacketOut, packetBytes, 0o644); err != nil { return err }

	decodedPayload, _, err := selfcontainedbook.DecodePNG(rendered.PNG); if err != nil { return err }
	decodedHeader, decodedCompressed, err := selfcontainedbook.Unpack(decodedPayload); if err != nil { return err }
	recovered, err := bzip2Decompress(decodedCompressed); if err != nil { return err }
	exact := bytes.Equal(text, recovered) && sourceHash(recovered) == decodedHeader.SourceSHA256 && decodedHeader.MasterPromptSHA256 == promptSHA
	status := "PASS_ROUNDTRIP_PASS_512K"; if rendered.OverTargetBytes > 0 { status = "PASS_ROUNDTRIP_FAIL_512K" }; if !exact { status = "FAIL_ROUNDTRIP_OR_BINDING" }
	rep := report{Schema: selfcontainedbook.Schema + ".report", Profile: selfcontainedbook.ProfileID, Status: status, SourcePath: path, SourcePlane: "PDF_EXTRACTED_TEXT_LAYOUT_PLANE", Header: header, Control: control, PNGBytes: rendered.PNGBytes, TargetPNGBytes: rendered.TargetPNGBytes, OverTargetBytes: rendered.OverTargetBytes, DataWidthBits: rendered.DataWidthBits, DataHeightBits: rendered.DataHeightBits, RoundtripExact: exact, SemanticFoldClaim: "Fixed Carrier R2 BOOT/ROSETTA/root index + canonical page/block address graph", ExactResidualClaim: "lossless bzip2 residual encoded in pixels; literal transport is not semantic Fold", MasterPromptPath: masterPromptOut, MasterPromptSHA256: promptSHA, ModelPacketPath: modelPacketOut}
	return printJSON(rep)
}

func decode(path, out string) error {
	data, err := os.ReadFile(path); if err != nil { return err }
	payload, controlImage, err := selfcontainedbook.DecodePNG(data); if err != nil { return err }
	var controlBuf bytes.Buffer; if err := pngEncode(&controlBuf, controlImage); err != nil { return err }
	control, err := fixedcarrier.DecodePNG(controlBuf.Bytes()); if err != nil { return err }
	header, compressed, err := selfcontainedbook.Unpack(payload); if err != nil { return err }
	recovered, err := bzip2Decompress(compressed); if err != nil { return err }
	if uint64(len(recovered)) != header.SourceBytes || sourceHash(recovered) != header.SourceSHA256 { return fmt.Errorf("recovered source verification failed") }
	if err := os.WriteFile(out, recovered, 0o644); err != nil { return err }
	return printJSON(struct { Schema string `json:"schema"`; Header selfcontainedbook.Header `json:"header"`; Control fixedcarrier.Decoded `json:"control"`; Output string `json:"output"`; Exact bool `json:"exact"` }{selfcontainedbook.Schema + ".decode", header, control, out, true})
}

func extractPDFText(path string) ([]byte, error) {
	binary, err := exec.LookPath("pdftotext"); if err != nil { return nil, fmt.Errorf("PDF text extraction requires pdftotext: %w", err) }
	cmd := exec.Command(binary, "-layout", "-enc", "UTF-8", path, "-"); body, err := cmd.Output(); if err != nil { return nil, fmt.Errorf("pdftotext failed: %w", err) }
	if strings.TrimSpace(string(body)) == "" { return nil, fmt.Errorf("PDF has no extractable text") }
	return body, nil
}

func bzip2Compress(data []byte) ([]byte, error) {
	binary, err := exec.LookPath("bzip2"); if err != nil { return nil, fmt.Errorf("self-contained book R0 requires bzip2 for the exact residual baseline: %w", err) }
	cmd := exec.Command(binary, "-9", "-c"); cmd.Stdin = bytes.NewReader(data); out, err := cmd.Output(); if err != nil { return nil, fmt.Errorf("bzip2 compression failed: %w", err) }; return out, nil
}
func bzip2Decompress(data []byte) ([]byte, error) {
	binary, err := exec.LookPath("bzip2"); if err != nil { return nil, fmt.Errorf("self-contained book R0 requires bzip2: %w", err) }
	cmd := exec.Command(binary, "-d", "-c"); cmd.Stdin = bytes.NewReader(data); out, err := cmd.Output(); if err != nil { return nil, fmt.Errorf("bzip2 decompression failed: %w", err) }; return out, nil
}
func sourceHash(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func hashBytes(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func printJSON(value any) error { body, err := json.MarshalIndent(value, "", "  "); if err != nil { return err }; fmt.Println(string(body)); return nil }
func pngEncode(buf *bytes.Buffer, img image.Image) error { enc := png.Encoder{CompressionLevel: png.BestCompression}; return enc.Encode(buf, img) }
func die(err error) { if err != nil { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) } }
