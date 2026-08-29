package spec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const SchemaV1 = "ohf.experiment.v1"

type ExperimentSpec struct {
	Schema       string       `json:"schema"`
	ID           string       `json:"id"`
	Baseline     string       `json:"baseline"`
	Campaign     string       `json:"campaign"`
	Seed         SeedSpec     `json:"seed"`
	Source       SourceSpec   `json:"source"`
	Renderer     RendererSpec `json:"renderer"`
	Variants     int          `json:"variants"`
	Degradations []string     `json:"degradations"`
	Evaluators   []string     `json:"evaluators"`
}
type SeedSpec struct {
	Algorithm string `json:"algorithm"`
	Value     uint64 `json:"value"`
}
type SourceSpec struct {
	Kind      string `json:"kind"`
	Generator string `json:"generator"`
	Length    int    `json:"length"`
	Path      string `json:"path,omitempty"`
	PageStart int    `json:"page_start,omitempty"`
	PageEnd   int    `json:"page_end,omitempty"`
}
type RendererSpec struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

func Load(path string) (ExperimentSpec, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return ExperimentSpec{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var s ExperimentSpec
	if err := dec.Decode(&s); err != nil {
		return ExperimentSpec{}, fmt.Errorf("decode spec: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		if err == nil {
			return ExperimentSpec{}, errors.New("decode spec: trailing JSON content")
		}
		return ExperimentSpec{}, fmt.Errorf("decode spec: trailing content: %w", err)
	}
	if err := s.Validate(); err != nil {
		return ExperimentSpec{}, err
	}
	return s, nil
}
func (s ExperimentSpec) Validate() error {
	if s.Schema != SchemaV1 {
		return fmt.Errorf("schema: expected %q, got %q", SchemaV1, s.Schema)
	}
	if s.ID == "" {
		return errors.New("id: required")
	}
	if s.Baseline != "R3.10" {
		return fmt.Errorf("baseline: F1 accepts only R3.10, got %q", s.Baseline)
	}
	if s.Campaign == "" {
		return errors.New("campaign: required")
	}
	if s.Seed.Algorithm != "splitmix64-v1" {
		return fmt.Errorf("seed.algorithm: unsupported %q", s.Seed.Algorithm)
	}
	switch s.Source.Kind {
	case "synthetic":
		if s.Source.Generator != "fixed-bytes-v1" {
			return fmt.Errorf("source.generator: synthetic source requires fixed-bytes-v1, got %q", s.Source.Generator)
		}
		if s.Source.Length <= 0 {
			return errors.New("source.length: must be > 0 for synthetic source")
		}
		if s.Source.Path != "" {
			return errors.New("source.path: must be empty for synthetic source")
		}
		if s.Source.PageStart != 0 || s.Source.PageEnd != 0 {
			return errors.New("source.page_start and source.page_end: must be empty for synthetic source")
		}
	case "file":
		if s.Source.Generator != "pdf-extract-v1" {
			return fmt.Errorf("source.generator: file source requires pdf-extract-v1, got %q", s.Source.Generator)
		}
		if s.Source.Path == "" {
			return errors.New("source.path: required for file source")
		}
		if s.Source.Length != 0 {
			return errors.New("source.length: must be 0 for file source")
		}
		if s.Source.PageStart < 0 || s.Source.PageEnd < 0 || (s.Source.PageEnd > 0 && s.Source.PageEnd < max(1, s.Source.PageStart)) {
			return errors.New("source page range: expected 1 <= page_start <= page_end")
		}
	default:
		return fmt.Errorf("source.kind: unsupported %q", s.Source.Kind)
	}
	if s.Renderer.ID == "" || s.Renderer.Version == "" {
		return errors.New("renderer.id and renderer.version: required")
	}
	if s.Renderer.Width <= 0 || s.Renderer.Height <= 0 {
		return errors.New("renderer.width and renderer.height: must be > 0")
	}
	if s.Variants <= 0 {
		return errors.New("variants: must be > 0")
	}
	if len(s.Evaluators) == 0 {
		return errors.New("evaluators: at least one evaluator required")
	}
	for _, e := range s.Evaluators {
		if e != "computational" && e != "native" {
			return fmt.Errorf("evaluators: unsupported %q", e)
		}
	}
	return nil
}
