package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/codec"
)

func main() {
	mode := flag.String("mode", "registry", "registry|encode-index|decode-index|roundtrip-index")
	in := flag.String("in", "", "input JSON for index/construction IR modes")
	out := flag.String("out", "-", "output JSON path or - for stdout")
	profile := flag.String("profile", "origami.fixed-carrier.r2.profile-3", "profile ID for encoding")
	flag.Parse()

	var value any
	switch *mode {
	case "registry":
		r := codec.ReferenceRegistry()
		die(r.Validate())
		value = r
	case "encode-index":
		var idx codec.SemanticIndex
		readJSON(*in, &idx)
		ir, err := codec.EncodeSuperIndex(idx, *profile)
		die(err)
		value = ir
	case "decode-index":
		var ir codec.ConstructionIR
		readJSON(*in, &ir)
		idx, err := codec.DecodeSuperIndex(ir)
		die(err)
		value = idx
	case "roundtrip-index":
		var idx codec.SemanticIndex
		readJSON(*in, &idx)
		ir, err := codec.EncodeSuperIndex(idx, *profile)
		die(err)
		decoded, err := codec.DecodeSuperIndex(ir)
		die(err)
		if !codec.EqualSemanticIndex(idx, decoded) {
			die(fmt.Errorf("S2/E2 semantic roundtrip mismatch"))
		}
		value = struct {
			Status       string               `json:"status"`
			Equation     string               `json:"equation"`
			Construction codec.ConstructionIR `json:"construction_ir"`
			Decoded      codec.SemanticIndex   `json:"decoded"`
		}{Status: "PASS", Equation: "S2(E2(INDEX)) ~= INDEX", Construction: ir, Decoded: decoded}
	default:
		die(fmt.Errorf("unknown mode %q", *mode))
	}
	writeJSON(*out, value)
}

func readJSON(path string, dst any) {
	if path == "" {
		die(fmt.Errorf("-in required"))
	}
	b, err := os.ReadFile(path)
	die(err)
	die(json.Unmarshal(b, dst))
}

func writeJSON(path string, value any) {
	b, err := json.MarshalIndent(value, "", "  ")
	die(err)
	b = append(b, '\n')
	if path == "-" {
		_, err = os.Stdout.Write(b)
		die(err)
		return
	}
	die(os.WriteFile(path, b, 0644))
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
