package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/receiver"
)

func main() {
	in := flag.String("in", "", "Tlaloc hybrid artifact-set JSON")
	root := flag.String("root", "receiver/candidates", "Origami candidate registry root")
	flag.Parse()
	if *in == "" {
		fmt.Fprintln(os.Stderr, "-in is required")
		os.Exit(2)
	}
	b, err := os.ReadFile(*in)
	must(err)
	var proposal receiver.TlalocHybridArtifactSet
	must(json.Unmarshal(b, &proposal))
	candidate, err := receiver.ImportTlalocCandidate(proposal)
	must(err)
	stored, err := receiver.StoreImportedCandidate(*root, candidate)
	must(err)
	fmt.Printf("CANDIDATE_ID=%s\n", candidate.CandidateID)
	fmt.Printf("STATUS=%s\n", receiver.ArtifactReferenceCandidate)
	fmt.Printf("PROMPT=%s\n", stored.PromptPath)
	fmt.Printf("CANDIDATE=%s\n", stored.CandidatePath)
	fmt.Printf("MANIFEST=%s\n", stored.ManifestPath)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
