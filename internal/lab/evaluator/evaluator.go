package evaluator

import "fmt"

type Mode string

const (
	Computational Mode = "computational"
	Native        Mode = "native"
)

type Capabilities struct {
	GroundTruth  bool
	PixelAccess  bool
	RunMetadata  bool
	GoDecoder    bool
	MasterPrompt bool
	CarrierImage bool
	Query        bool
}

func For(mode Mode) (Capabilities, error) {
	switch mode {
	case Computational:
		return Capabilities{GroundTruth: true, PixelAccess: true, RunMetadata: true, GoDecoder: true}, nil
	case Native:
		return Capabilities{MasterPrompt: true, CarrierImage: true, Query: true}, nil
	default:
		return Capabilities{}, fmt.Errorf("unknown evaluator mode %q", mode)
	}
}

func ValidateIsolation() error {
	n, _ := For(Native)
	if n.GroundTruth || n.PixelAccess || n.RunMetadata || n.GoDecoder {
		return fmt.Errorf("native evaluator isolation violated")
	}
	if !n.MasterPrompt || !n.CarrierImage || !n.Query {
		return fmt.Errorf("native evaluator missing required blind inputs")
	}
	return nil
}
