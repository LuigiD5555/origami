package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/temporal"
)

func main() {
	mode := flag.String("mode", "run", "run|seek|verify|metrics")
	in := flag.String("in", "", "temporal program JSON")
	tracePath := flag.String("trace", "", "trace JSON for seek/verify/metrics")
	step := flag.Int("step", 0, "target step for seek")
	out := flag.String("out", "-", "output JSON path or -")
	flag.Parse()
	if *in == "" && *mode != "metrics" { fatal("-in is required") }

	var program temporal.Program
	if *in != "" { readJSON(*in, &program) }

	switch *mode {
	case "run":
		trace, err := temporal.Run(program); if err != nil { fatal(err.Error()) }
		writeJSON(*out, trace)
	case "seek":
		if *tracePath == "" { fatal("-trace is required") }
		var trace temporal.Trace; readJSON(*tracePath, &trace)
		state, err := temporal.Seek(trace, *step); if err != nil { fatal(err.Error()) }
		writeJSON(*out, map[string]any{"step": *step, "state": state})
	case "verify":
		if *tracePath == "" { fatal("-trace is required") }
		var trace temporal.Trace; readJSON(*tracePath, &trace)
		if err := temporal.VerifyReplay(program, trace); err != nil { fatal(err.Error()) }
		writeJSON(*out, map[string]any{"status": "PASS_REPLAY_EQUIVALENCE", "program_id": program.ID})
	case "metrics":
		if *tracePath == "" { fatal("-trace is required") }
		var trace temporal.Trace; readJSON(*tracePath, &trace)
		writeJSON(*out, temporal.Measure(trace))
	default:
		fatal("unsupported -mode")
	}
}

func readJSON(path string, out any) {
	b, err := os.ReadFile(path); if err != nil { fatal(err.Error()) }
	if err := json.Unmarshal(b, out); err != nil { fatal(err.Error()) }
}

func writeJSON(path string, v any) {
	var f *os.File
	if path == "" || path == "-" { f = os.Stdout } else {
		var err error; f, err = os.Create(path); if err != nil { fatal(err.Error()) }; defer f.Close()
	}
	enc := json.NewEncoder(f); enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil { fatal(err.Error()) }
}

func fatal(msg string) { fmt.Fprintln(os.Stderr, msg); os.Exit(2) }
