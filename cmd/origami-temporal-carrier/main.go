package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func main() {
	mode := flag.String("mode", "build", "build|decode")
	in := flag.String("in", "", "temporal program JSON or PNG")
	out := flag.String("out", "", "output PNG or recovered program JSON")
	flag.Parse()
	if *in == "" || *out == "" { fatal("-in and -out are required") }

	switch *mode {
	case "build":
		var p temporal.Program
		readJSON(*in, &p)
		pngData, info, err := fixedcarrier.RenderTemporalCarrier(p)
		if err != nil { fatal(err.Error()) }
		if err := os.WriteFile(*out, pngData, 0644); err != nil { fatal(err.Error()) }
		writeStdout(info)
	case "decode":
		b, err := os.ReadFile(*in); if err != nil { fatal(err.Error()) }
		info, err := fixedcarrier.DecodeTemporalCarrierPNG(b); if err != nil { fatal(err.Error()) }
		writeJSON(*out, info.Program)
		writeStdout(info)
	default:
		fatal("unsupported -mode")
	}
}

func readJSON(path string, out any) {
	b, err := os.ReadFile(path); if err != nil { fatal(err.Error()) }
	if err := json.Unmarshal(b, out); err != nil { fatal(err.Error()) }
}
func writeJSON(path string, v any) {
	b, err := json.MarshalIndent(v, "", "  "); if err != nil { fatal(err.Error()) }
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0644); err != nil { fatal(err.Error()) }
}
func writeStdout(v any) {
	enc := json.NewEncoder(os.Stdout); enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil { fatal(err.Error()) }
}
func fatal(msg string) { fmt.Fprintln(os.Stderr, msg); os.Exit(2) }
