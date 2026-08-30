package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LuigiD5555/origami/internal/evidence"
	"github.com/LuigiD5555/origami/internal/memory"
)

type input struct {
	Schema    string              `json:"schema,omitempty"`
	CarrierID string              `json:"carrier_id"`
	Nodes     []memory.Node       `json:"nodes"`
	Policy    evidence.Policy     `json:"policy,omitempty"`
	Proposals []evidence.Proposal `json:"proposals"`
}

func main() {
	inPath := flag.String("in", "-", "input JSON path or - for stdin")
	outPath := flag.String("out", "-", "output JSON path or - for stdout")
	flag.Parse()

	var r io.Reader = os.Stdin
	if *inPath != "-" {
		f, err := os.Open(*inPath)
		die(err)
		defer f.Close()
		r = f
	}
	var payload input
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	die(dec.Decode(&payload))
	if payload.CarrierID == "" {
		payload.CarrierID = "canonical"
	}
	store, err := memory.BuildStore(payload.CarrierID, payload.Nodes)
	die(err)
	result, err := evidence.NewReducer(evidence.NewMemoryResolver(store), payload.Policy).Reduce(payload.Proposals)
	die(err)
	body, err := json.MarshalIndent(result, "", "  ")
	die(err)
	body = append(body, '\n')
	if *outPath == "-" {
		_, err = os.Stdout.Write(body)
		die(err)
		return
	}
	die(os.WriteFile(*outPath, body, 0o644))
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
