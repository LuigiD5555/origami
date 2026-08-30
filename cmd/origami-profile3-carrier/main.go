package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

func main() {
	mode := flag.String("mode", "build", "build|decode")
	in := flag.String("in", "", "metadata JSON for build or carrier PNG for decode")
	out := flag.String("out", "origami-profile3.png", "output PNG for build")
	flag.Parse()

	switch *mode {
	case "build":
		if *in == "" { die(fmt.Errorf("-in metadata json required")) }
		b, err := os.ReadFile(*in); die(err)
		var meta fixedcarrier.Metadata; die(json.Unmarshal(b, &meta))
		dec, n, err := fixedcarrier.WriteProfile3PNG(*out, meta); die(err)
		jb, _ := json.MarshalIndent(dec, "", "  ")
		fmt.Printf("%s\nPNG_BYTES=%d\n", jb, n)
	case "decode":
		if *in == "" { die(fmt.Errorf("-in carrier png required")) }
		b, err := os.ReadFile(*in); die(err)
		dec, err := fixedcarrier.DecodeProfile3PNG(b); die(err)
		jb, _ := json.MarshalIndent(dec, "", "  ")
		fmt.Println(string(jb))
	default:
		die(fmt.Errorf("unknown mode %q", *mode))
	}
}

func die(err error) { if err != nil { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) } }
