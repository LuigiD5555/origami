package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LuigiD5555/origami/internal/document"
)

func main() {
	in := flag.String("in", "", "source path")
	out := flag.String("out", "-", "bundle JSON path or - for stdout")
	mode := flag.String("mode", "auto", "auto|text|pdf|layout-json")
	sourceID := flag.String("source-id", "", "stable source id; defaults to input basename")
	mediaType := flag.String("media-type", "", "media type override")
	flag.Parse()
	if *in == "" { die(fmt.Errorf("-in is required")) }
	id := *sourceID
	if id == "" { id = filepath.Base(*in) }

	resolvedMode := *mode
	if resolvedMode == "auto" {
		switch strings.ToLower(filepath.Ext(*in)) {
		case ".pdf": resolvedMode = "pdf"
		case ".json": resolvedMode = "layout-json"
		default: resolvedMode = "text"
		}
	}

	var doc document.Document
	var err error
	switch resolvedMode {
	case "text":
		body, readErr := os.ReadFile(*in); die(readErr)
		mt := *mediaType; if mt == "" { mt = "text/plain" }
		doc, err = document.FromText(id, mt, body)
	case "pdf":
		body, extractErr := extractPDFText(*in); die(extractErr)
		mt := *mediaType; if mt == "" { mt = "application/pdf" }
		doc, err = document.FromText(id, mt, body)
	case "layout-json":
		body, readErr := os.ReadFile(*in); die(readErr)
		doc, err = document.LoadJSON(body)
	default:
		die(fmt.Errorf("unsupported mode %q", resolvedMode))
	}
	die(err)
	bundle, err := document.BuildBundle(doc); die(err)
	body, err := json.MarshalIndent(bundle, "", "  "); die(err)
	body = append(body, '\n')
	if *out == "-" {
		_, err = os.Stdout.Write(body); die(err)
		return
	}
	die(os.WriteFile(*out, body, 0o644))
}

func extractPDFText(path string) ([]byte, error) {
	binary, err := exec.LookPath("pdftotext")
	if err != nil {
		return nil, fmt.Errorf("PDF text extraction requires pdftotext; no screenshot/raster fallback is allowed: %w", err)
	}
	cmd := exec.Command(binary, "-layout", "-enc", "UTF-8", path, "-")
	body, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pdftotext failed: %w", err)
	}
	if strings.TrimSpace(string(body)) == "" {
		return nil, fmt.Errorf("PDF has no extractable text; use an OCR/layout adapter that emits origami.document.r0 instead of rasterizing the page as semantic content")
	}
	return body, nil
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
