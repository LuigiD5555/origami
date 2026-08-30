package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/writer"
)

func main() {
	in := flag.String("in", "-", "writer request JSON path or - for stdin")
	out := flag.String("out", "-", "construction plan JSON path or - for stdout")
	flag.Parse()
	body, err := readInput(*in)
	die(err)
	var request writer.Request
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	die(dec.Decode(&request))
	plan, err := writer.Plan(request)
	die(err)
	encoded, err := json.MarshalIndent(plan, "", "  ")
	die(err)
	encoded = append(encoded, '\n')
	if *out == "-" {
		_, err = os.Stdout.Write(encoded)
		die(err)
		return
	}
	die(os.WriteFile(*out, encoded, 0o644))
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return os.ReadFile("/dev/stdin")
	}
	return os.ReadFile(path)
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
