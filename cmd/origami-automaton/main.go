package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/automaton"
)

func main() {
	mode := flag.String("mode", "inspect", "inspect|step")
	in := flag.String("in", "", "automaton JSON")
	statePath := flag.String("state", "", "optional state JSON for step")
	flag.Parse()
	if *in == "" { fatal("-in is required") }

	var a automaton.Automaton
	readJSON(*in, &a)
	if err := a.Validate(); err != nil { fatal(err.Error()) }

	switch *mode {
	case "inspect":
		writeJSON(map[string]any{"schema": a.Schema, "id": a.ID, "cells": len(a.Cells), "rules": len(a.Rules), "edges": len(a.Edges), "initial_state": a.Initial()})
	case "step":
		state := a.Initial()
		if *statePath != "" { readJSON(*statePath, &state) }
		next, result, err := automaton.Step(a, state, 1)
		if err != nil { fatal(err.Error()) }
		writeJSON(map[string]any{"result": result, "state": next})
	default:
		fatal("unsupported -mode")
	}
}

func readJSON(path string, out any) {
	b, err := os.ReadFile(path); if err != nil { fatal(err.Error()) }
	if err := json.Unmarshal(b, out); err != nil { fatal(err.Error()) }
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout); enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil { fatal(err.Error()) }
}

func fatal(msg string) { fmt.Fprintln(os.Stderr, msg); os.Exit(2) }
