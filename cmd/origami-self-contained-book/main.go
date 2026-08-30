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
	Schema                     string                   `json:"schema"`
	Profile                    string                   `json:"profile"`
	Status                     string                   `json:"status"`
	SourcePath                 string                   `json:"source_path,omitempty"`
	SourcePlane                string                   `json:"source_plane"`
	Header                     selfcontainedbook.Header `json:"header"`
	Control                    fixedcarrier.Decoded     `json:"control"`
	PNGBytes                   int                      `json:"png_bytes"`
	TargetPNGBytes             int                      `json:"target_png_bytes"`
	OverTargetBytes            int                      `json:"over_target_bytes"`
	DataWidthBits              int                      `json:"data_width_bits"`
	DataHeightBits             int                      `json:"data_height_bits"`
	RoundtripExact             bool                     `json:"roundtrip_exact"`
	MasterPromptRoundtripExact bool                     `json:"master_prompt_roundtrip_exact"`
	MasterPromptEmbedded       bool                     `json:"master_prompt_embedded"`
	ExternalPromptRequired     bool                     `json:"external_prompt_required"`
	SemanticFoldClaim          string                   `json:"semantic_fold_claim"`
	ExactResidualClaim         string                   `json:"exact_residual_claim"`
	MasterPromptExport         string                   `json:"master_prompt_export"`
	MasterPromptSHA256         string                   `json:"master_prompt_sha256"`
	ModelPacketPath            string                   `json:"model_packet_path"`
}

type modelPacket struct {
	Schema                 string   `json:"schema"`
	Profile                string   `json:"profile"`
	Carrier                string   `json:"carrier"`
	CarrierSHA256          string   `json:"carrier_sha256"`
	MasterPromptEmbedded   bool     `json:"master_prompt_embedded"`
	ExternalPromptRequired bool     `json:"external_prompt_required"`
	MasterPromptExport     string   `json:"master_prompt_export"`
	MasterPromptSHA256     string   `json:"master_prompt_sha256"`
	SourceSHA256           string   `json:"source_sha256"`
	WorkingWindowTokenEq   int      `json:"working_window_token_eq"`
	RequiredInputs         []string `json:"required_inputs"`
	ReadOrder              []string `json:"read_order"`
	Claims                 []string `json:"claims"`
}

func main() {
	mode := flag.String("mode", "build", "build|decode|prompt")
	in := flag.String("in", "", "PDF input for build or carrier PNG for decode/prompt")
	out := flag.String("out", "origami-self-contained-book-r0.png", "carrier PNG for build, recovered text for decode, or prompt path/- for prompt")
	prompt := flag.String("prompt", "generated/MASTER_PROMPT.md", "Master Prompt source used while building the self-contained image")
	promptOut := flag.String("prompt-out", "", "decode mode: recovered Master Prompt path; defaults beside recovered text")
	dataWidth := flag.Int("data-width", selfcontainedbook.DefaultDataWidth, "exact payload matrix width in bits")
	flag.Parse()
	if *in == "" {
		die(fmt.Errorf("-in is required"))
	}
	switch *mode {
	case "build":
		die(build(*in, *out, *prompt, *dataWidth))
	case "decode":
		die(decode(*in, *out, *promptOut))
	case "prompt":
		die(extractPrompt(*in, *out))
	default:
		die(fmt.Errorf("unsupported mode %q", *mode))
	}
}

func build(path, out, promptPath string, dataWidth int) error {
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("read Master Prompt: %w", err)
	}
	promptSHA := hashBytes(promptBytes)
	text, err := extractPDFText(path)
	if err != nil {
		return err
	}
	doc, err := document.FromText(filepath.Base(path), "application/pdf", text)
	if err != nil {
		return err
	}
	bundle, err := document.BuildBundle(doc)
	if err != nil {
		return err
	}
	system, err := integration.NewSystem(doc.SourceID, bundle.Nodes)
	if err != nil {
		return err
	}
	blockCount := 0
	for _, page := range doc.Pages {
		blockCount += len(page.Blocks)
	}
	signature, err := json.Marshal(system.Signature())
	if err != nil {
		return err
	}
	compressed, err := bzip2Compress(text)
	if err != nil {
		return err
	}
	payload, header, err := selfcontainedbook.Pack(text, promptBytes, compressed, uint32(len(doc.Pages)), uint32(blockCount), "bz2-9")
	if err != nil {
		return err
	}
	meta := fixedcarrier.Metadata{
		CarrierID:      doc.SourceHash,
		StoreRoot:      integration.CanonicalStoreRoot(system.Nodes()),
		SourceSHA256:   doc.SourceHash,
		PageCount:      uint32(len(doc.Pages)),
		BlockCount:     uint32(blockCount),
		DocumentCount:  1,
		ObjectCount:    uint32(len(bundle.Nodes)),
		GraphSignature: signature,
	}
	controlPNG, control, err := fixedcarrier.Render(meta)
	if err != nil {
		return err
	}
	rendered, err := selfcontainedbook.Render(controlPNG, payload, dataWidth)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(out, rendered.PNG, 0o644); err != nil {
		return err
	}

	// Prove that both the Master Prompt and book plane come back from the PNG,
	// then export the decoded prompt only as a convenience/debug artifact.
	decodedPayload, _, err := selfcontainedbook.DecodePNG(rendered.PNG)
	if err != nil {
		return err
	}
	decodedHeader, decodedPrompt, decodedCompressed, err := selfcontainedbook.Unpack(decodedPayload)
	if err != nil {
		return err
	}
	recovered, err := bzip2Decompress(decodedCompressed)
	if err != nil {
		return err
	}
	sourceExact := bytes.Equal(text, recovered) && sourceHash(recovered) == decodedHeader.SourceSHA256
	promptExact := bytes.Equal(promptBytes, decodedPrompt) && hashBytes(decodedPrompt) == decodedHeader.MasterPromptSHA256

	packageDir := filepath.Dir(out)
	masterPromptExport := filepath.Join(packageDir, "MASTER_PROMPT.extracted.md")
	modelPacketOut := filepath.Join(packageDir, "model_packet.json")
	if err := os.WriteFile(masterPromptExport, decodedPrompt, 0o644); err != nil {
		return err
	}
	packet := modelPacket{
		Schema:                 "origami.self-contained-book.r0-lab.model-packet",
		Profile:                selfcontainedbook.ProfileID,
		Carrier:                filepath.Base(out),
		CarrierSHA256:          hashBytes(rendered.PNG),
		MasterPromptEmbedded:   true,
		ExternalPromptRequired: false,
		MasterPromptExport:     filepath.Base(masterPromptExport),
		MasterPromptSHA256:     decodedHeader.MasterPromptSHA256,
		SourceSHA256:           decodedHeader.SourceSHA256,
		WorkingWindowTokenEq:   4000,
		RequiredInputs:         []string{filepath.Base(out), "user question"},
		ReadOrder:              []string{"S0_SELF_BOOT", "MP0_MASTER_PROMPT", "T0_BOOT", "T1_ROSETTA", "T2_SUPERINDEX", "T3_MEMORY", "VERIFY"},
		Claims: []string{
			"MASTER_PROMPT_BYTES_ARE_INSIDE_CARRIER",
			"MASTER_PROMPT_HASH_VERIFIED",
			"EXTERNAL_MASTER_PROMPT_NOT_REQUIRED_FOR_MECHANICAL_DECODE",
			"NATIVE_MODEL_ONLY_BOOT_REMAINS_EMPIRICAL_NOT_PROMOTED",
			"LITERAL_RESIDUAL_IS_NOT_SEMANTIC_FOLD",
		},
	}
	packetBytes, err := json.MarshalIndent(packet, "", "  ")
	if err != nil {
		return err
	}
	packetBytes = append(packetBytes, '\n')
	if err := os.WriteFile(modelPacketOut, packetBytes, 0o644); err != nil {
		return err
	}

	exact := sourceExact && promptExact
	status := "PASS_SELF_CONTAINED_ROUNDTRIP_PASS_512K"
	if rendered.OverTargetBytes > 0 {
		status = "PASS_SELF_CONTAINED_ROUNDTRIP_FAIL_512K"
	}
	if !exact {
		status = "FAIL_SELF_CONTAINED_ROUNDTRIP"
	}
	rep := report{
		Schema:                     selfcontainedbook.Schema + ".report",
		Profile:                    selfcontainedbook.ProfileID,
		Status:                     status,
		SourcePath:                 path,
		SourcePlane:                "PDF_EXTRACTED_TEXT_LAYOUT_PLANE",
		Header:                     header,
		Control:                    control,
		PNGBytes:                   rendered.PNGBytes,
		TargetPNGBytes:             rendered.TargetPNGBytes,
		OverTargetBytes:            rendered.OverTargetBytes,
		DataWidthBits:              rendered.DataWidthBits,
		DataHeightBits:             rendered.DataHeightBits,
		RoundtripExact:             sourceExact,
		MasterPromptRoundtripExact: promptExact,
		MasterPromptEmbedded:       true,
		ExternalPromptRequired:     false,
		SemanticFoldClaim:          "Fixed Carrier R2 BOOT/ROSETTA/root index + canonical page/block address graph; visible S0/MP0 self-bootstrap wrapper",
		ExactResidualClaim:         "book body is still a lossless bzip2 residual encoded in pixels; literal transport is not semantic Fold",
		MasterPromptExport:         masterPromptExport,
		MasterPromptSHA256:         promptSHA,
		ModelPacketPath:            modelPacketOut,
	}
	return printJSON(rep)
}

func decode(path, out, promptOut string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	payload, controlImage, err := selfcontainedbook.DecodePNG(data)
	if err != nil {
		return err
	}
	var controlBuf bytes.Buffer
	if err := pngEncode(&controlBuf, controlImage); err != nil {
		return err
	}
	control, err := fixedcarrier.DecodePNG(controlBuf.Bytes())
	if err != nil {
		return err
	}
	header, masterPrompt, compressed, err := selfcontainedbook.Unpack(payload)
	if err != nil {
		return err
	}
	recovered, err := bzip2Decompress(compressed)
	if err != nil {
		return err
	}
	if uint64(len(recovered)) != header.SourceBytes || sourceHash(recovered) != header.SourceSHA256 {
		return fmt.Errorf("recovered source verification failed")
	}
	if uint32(len(masterPrompt)) != header.MasterPromptBytes || hashBytes(masterPrompt) != header.MasterPromptSHA256 {
		return fmt.Errorf("recovered Master Prompt verification failed")
	}
	if err := os.WriteFile(out, recovered, 0o644); err != nil {
		return err
	}
	if promptOut == "" {
		promptOut = filepath.Join(filepath.Dir(out), "MASTER_PROMPT.recovered.md")
	}
	if err := os.WriteFile(promptOut, masterPrompt, 0o644); err != nil {
		return err
	}
	return printJSON(struct {
		Schema       string                   `json:"schema"`
		Header       selfcontainedbook.Header `json:"header"`
		Control      fixedcarrier.Decoded     `json:"control"`
		Output       string                   `json:"output"`
		PromptOutput string                   `json:"prompt_output"`
		SourceExact  bool                     `json:"source_exact"`
		PromptExact  bool                     `json:"prompt_exact"`
	}{selfcontainedbook.Schema + ".decode", header, control, out, promptOut, true, true})
}

func extractPrompt(path, out string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	payload, _, err := selfcontainedbook.DecodePNG(data)
	if err != nil {
		return err
	}
	header, prompt, _, err := selfcontainedbook.Unpack(payload)
	if err != nil {
		return err
	}
	if hashBytes(prompt) != header.MasterPromptSHA256 {
		return fmt.Errorf("Master Prompt verification failed")
	}
	if out == "-" {
		_, err = os.Stdout.Write(prompt)
		return err
	}
	return os.WriteFile(out, prompt, 0o644)
}

func extractPDFText(path string) ([]byte, error) {
	binary, err := exec.LookPath("pdftotext")
	if err != nil {
		return nil, fmt.Errorf("PDF text extraction requires pdftotext: %w", err)
	}
	cmd := exec.Command(binary, "-layout", "-enc", "UTF-8", path, "-")
	body, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pdftotext failed: %w", err)
	}
	if strings.TrimSpace(string(body)) == "" {
		return nil, fmt.Errorf("PDF has no extractable text")
	}
	return body, nil
}

func bzip2Compress(data []byte) ([]byte, error) {
	binary, err := exec.LookPath("bzip2")
	if err != nil {
		return nil, fmt.Errorf("self-contained book R0 requires bzip2 for the exact residual baseline: %w", err)
	}
	cmd := exec.Command(binary, "-9", "-c")
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("bzip2 compression failed: %w", err)
	}
	return out, nil
}

func bzip2Decompress(data []byte) ([]byte, error) {
	binary, err := exec.LookPath("bzip2")
	if err != nil {
		return nil, fmt.Errorf("self-contained book R0 requires bzip2: %w", err)
	}
	cmd := exec.Command(binary, "-d", "-c")
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("bzip2 decompression failed: %w", err)
	}
	return out, nil
}

func sourceHash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
func printJSON(value any) error {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
func pngEncode(buf *bytes.Buffer, img image.Image) error {
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	return enc.Encode(buf, img)
}
func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
