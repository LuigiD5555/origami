package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/lab/hybridcarrier"
	"github.com/LuigiD5555/origami/internal/receiver"
)

type Input struct {
	Receiver receiver.Spec               `json:"receiver"`
	Index    []hybridcarrier.IndexEntry  `json:"index"`
	Memory   []hybridcarrier.MemoryEntry `json:"memory"`
}

type Transport struct {
	Columns int `json:"columns"`
	Cell    int `json:"cell_pixels"`
	Margin  int `json:"margin"`
}

type ModelPacket struct {
	Schema               string    `json:"schema"`
	ReceiverContract     string    `json:"receiver_contract"`
	CarrierSchema        string    `json:"carrier_schema"`
	Mode                 string    `json:"mode"`
	WorkingWindowTokenEq int       `json:"working_window_token_eq"`
	MasterPromptSHA256   string    `json:"master_prompt_sha256"`
	CarrierSHA256        string    `json:"carrier_sha256"`
	Transport            Transport `json:"transport"`
	AllowedTools         []string  `json:"allowed_tools"`
	AntiContamination    bool      `json:"anti_contamination"`
}

func main() {
	in := flag.String("in", "experiments/hybrid-carrier-synthetic-r0.json", "Hybrid carrier logical input JSON")
	promptPath := flag.String("prompt", "generated/MASTER_PROMPT.md", "universal receiver bootstrap prompt")
	root := flag.String("root", "runs/hybrid-carrier-synthetic-r0", "prepared run directory")
	columns := flag.Int("columns", hybridcarrier.DefaultColumns, "glyph columns")
	cell := flag.Int("cell", hybridcarrier.DefaultCellPixels, "glyph cell pixels")
	margin := flag.Int("margin", hybridcarrier.DefaultMargin, "image margin pixels")
	maxBytes := flag.Int("max-bytes", hybridcarrier.DefaultMaxPNGBytes, "maximum PNG bytes")
	flag.Parse()

	inputBytes, err := os.ReadFile(*in)
	must(err)
	var input Input
	must(json.Unmarshal(inputBytes, &input))
	prompt, err := os.ReadFile(*promptPath)
	must(err)

	envelope, err := hybridcarrier.BuildEnvelope(input.Receiver, input.Index, input.Memory)
	must(err)
	rendered, err := hybridcarrier.Render(envelope, *columns, *cell, *margin, *maxBytes)
	must(err)

	publicDir := filepath.Join(*root, "public")
	privateDir := filepath.Join(*root, "private")
	must(os.MkdirAll(publicDir, 0o755))
	must(os.MkdirAll(privateDir, 0o755))

	carrierPath := filepath.Join(publicDir, "carrier.png")
	must(os.WriteFile(carrierPath, rendered.PNG, 0o644))
	masterPromptPath := filepath.Join(publicDir, "MASTER_PROMPT.md")
	must(os.WriteFile(masterPromptPath, prompt, 0o644))

	packet := ModelPacket{
		Schema:               "origami.model-packet.r0",
		ReceiverContract:     receiver.ContractID,
		CarrierSchema:        hybridcarrier.SchemaR0,
		Mode:                 "hybrid",
		WorkingWindowTokenEq: input.Receiver.Boot.WorkingWindow,
		MasterPromptSHA256:   hash(prompt),
		CarrierSHA256:        rendered.PNGSHA256,
		Transport:            Transport{Columns: *columns, Cell: *cell, Margin: *margin},
		AllowedTools: []string{
			"BOOT", "LOOKUP", "FOLLOW", "TRACE", "VERIFY", "STOP",
		},
		AntiContamination: true,
	}
	packetBytes, err := json.MarshalIndent(packet, "", "  ")
	must(err)
	packetBytes = append(packetBytes, '\n')
	must(os.WriteFile(filepath.Join(publicDir, "model_packet.json"), packetBytes, 0o644))

	metaBytes, err := json.MarshalIndent(rendered, "", "  ")
	must(err)
	metaBytes = append(metaBytes, '\n')
	must(os.WriteFile(filepath.Join(privateDir, "carrier.meta.json"), metaBytes, 0o644))

	envelopeBytes, err := json.MarshalIndent(envelope, "", "  ")
	must(err)
	envelopeBytes = append(envelopeBytes, '\n')
	must(os.WriteFile(filepath.Join(privateDir, "carrier.envelope.json"), envelopeBytes, 0o644))
	must(os.WriteFile(filepath.Join(privateDir, "source.json"), inputBytes, 0o644))

	// Verify the public image itself can reconstruct the logical carrier. This
	// check deliberately does not read private/carrier.envelope.json.
	decoded, err := hybridcarrier.DecodePNG(rendered.PNG, *columns, *cell, *margin)
	must(err)
	if decoded.Verification.MemorySHA256 != envelope.Verification.MemorySHA256 {
		must(fmt.Errorf("carrier PNG self-decode verification mismatch"))
	}

	fmt.Printf("PUBLIC_DIR=%s\n", publicDir)
	fmt.Printf("CARRIER=%s\n", carrierPath)
	fmt.Printf("MASTER_PROMPT=%s\n", masterPromptPath)
	fmt.Printf("MODEL_PACKET=%s\n", filepath.Join(publicDir, "model_packet.json"))
	fmt.Printf("PNG_BYTES=%d\n", rendered.PNGBytes)
	fmt.Printf("PNG_SHA256=%s\n", rendered.PNGSHA256)
	fmt.Printf("PAYLOAD_BYTES=%d\n", rendered.PayloadBytes)
	fmt.Printf("PAYLOAD_SHA256=%s\n", rendered.PayloadSHA256)
	fmt.Printf("PRIVATE_DEBUG=%s\n", privateDir)
}

func hash(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
