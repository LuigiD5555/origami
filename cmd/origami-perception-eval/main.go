package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/promotion"
)

func main() {
	carrier := flag.String("carrier", "", "canonical Fixed Carrier R2 PNG")
	observation := flag.String("observation", "", "model observation JSON")
	out := flag.String("out", "-", "report JSON path or - for stdout")
	flag.Parse()
	if *carrier == "" || *observation == "" {
		die(fmt.Errorf("-carrier and -observation are required"))
	}
	png, err := os.ReadFile(*carrier)
	die(err)
	body, err := os.ReadFile(*observation)
	die(err)
	var obs promotion.Observation
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	die(dec.Decode(&obs))
	report, err := promotion.Evaluate(png, obs)
	die(err)
	encoded, err := json.MarshalIndent(report, "", "  ")
	die(err)
	encoded = append(encoded, '\n')
	if *out == "-" {
		_, err = os.Stdout.Write(encoded)
		die(err)
		return
	}
	die(os.WriteFile(*out, encoded, 0o644))
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
