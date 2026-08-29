package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/lab/hybridcarrier"
	"github.com/LuigiD5555/origami/internal/receiver"
)

type Input struct {
	Receiver receiver.Spec                `json:"receiver"`
	Index    []hybridcarrier.IndexEntry   `json:"index"`
	Memory   []hybridcarrier.MemoryEntry  `json:"memory"`
}

func main() {
	in := flag.String("in", "experiments/hybrid-carrier-synthetic-r0.json", "Hybrid carrier logical input JSON")
	out := flag.String("out", "runs/hybrid-carrier-synthetic-r0/carrier.png", "output carrier PNG")
	columns := flag.Int("columns", 24, "glyph columns")
	cell := flag.Int("cell", 40, "glyph cell pixels")
	margin := flag.Int("margin", 8, "image margin pixels")
	maxBytes := flag.Int("max-bytes", hybridcarrier.DefaultMaxPNGBytes, "maximum PNG bytes")
	flag.Parse()

	b, err := os.ReadFile(*in)
	must(err)
	var input Input
	must(json.Unmarshal(b, &input))
	envelope, err := hybridcarrier.BuildEnvelope(input.Receiver, input.Index, input.Memory)
	must(err)
	rendered, err := hybridcarrier.Render(envelope, *columns, *cell, *margin, *maxBytes)
	must(err)

	must(os.MkdirAll(filepath.Dir(*out), 0o755))
	must(os.WriteFile(*out, rendered.PNG, 0o644))

	metaPath := *out + ".meta.json"
	meta, err := json.MarshalIndent(rendered, "", "  ")
	must(err)
	meta = append(meta, '\n')
	must(os.WriteFile(metaPath, meta, 0o644))

	envelopePath := *out + ".envelope.json"
	envelopeBytes, err := json.MarshalIndent(envelope, "", "  ")
	must(err)
	envelopeBytes = append(envelopeBytes, '\n')
	must(os.WriteFile(envelopePath, envelopeBytes, 0o644))

	fmt.Printf("CARRIER=%s\n", *out)
	fmt.Printf("PNG_BYTES=%d\n", rendered.PNGBytes)
	fmt.Printf("PNG_SHA256=%s\n", rendered.PNGSHA256)
	fmt.Printf("PAYLOAD_BYTES=%d\n", rendered.PayloadBytes)
	fmt.Printf("PAYLOAD_SHA256=%s\n", rendered.PayloadSHA256)
	fmt.Printf("META=%s\n", metaPath)
	fmt.Printf("ENVELOPE_DEBUG=%s\n", envelopePath)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
