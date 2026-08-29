package source

import (
	"encoding/binary"
	"fmt"
	"github.com/LuigiD5555/origami/internal/lab/seed"
	"github.com/LuigiD5555/origami/internal/lab/spec"
)

const GeneratorFixedBytesV1 = "fixed-bytes-v1"
const SeedNamespace = "source"

type Material struct {
	Bytes         []byte
	Generator     string
	SeedNamespace string
	SeedValue     uint64
	PageIndex     []PageEntry
}

func Generate(s spec.ExperimentSpec) (Material, error) {
	if s.Source.Kind == "file" && s.Source.Generator == GeneratorPDFExtractV1 {
		material, _, err := GeneratePDF(s)
		return material, err
	}
	if s.Source.Kind != "synthetic" {
		return Material{}, fmt.Errorf("source: unsupported kind %q", s.Source.Kind)
	}
	if s.Source.Generator != GeneratorFixedBytesV1 {
		return Material{}, fmt.Errorf("source: unsupported generator %q", s.Source.Generator)
	}
	if s.Source.Length <= 0 {
		return Material{}, fmt.Errorf("source: length must be > 0")
	}
	derived := seed.Derive(s.Seed.Value, SeedNamespace)
	rng := seed.New(derived)
	out := make([]byte, 0, s.Source.Length)
	var buf [8]byte
	for len(out) < s.Source.Length {
		binary.LittleEndian.PutUint64(buf[:], rng.Uint64())
		remaining := s.Source.Length - len(out)
		n := len(buf)
		if remaining < n {
			n = remaining
		}
		out = append(out, buf[:n]...)
	}
	return Material{Bytes: out, Generator: GeneratorFixedBytesV1, SeedNamespace: SeedNamespace, SeedValue: derived}, nil
}
