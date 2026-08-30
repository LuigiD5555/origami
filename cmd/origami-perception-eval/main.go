package main

import (
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
	dec := json.NewDecoder(bytesReader(body))
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

type byteReader struct { data []byte; off int }
func bytesReader(data []byte) *byteReader { return &byteReader{data:data} }
func (r *byteReader) Read(p []byte) (int,error) { if r.off>=len(r.data){return 0,os.ErrClosed}; n:=copy(p,r.data[r.off:]);r.off+=n;return n,nil }

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
