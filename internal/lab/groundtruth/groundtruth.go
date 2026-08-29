package groundtruth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	labsource "github.com/LuigiD5555/origami/internal/lab/source"
	"github.com/LuigiD5555/origami/internal/lab/spec"
)

const SchemaV1 = "ohf.ground-truth.v1"
const DefaultBlockSize = 64

type Block struct {
	Index  int    `json:"index"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	SHA256 string `json:"sha256"`
}

type SourceTruth struct {
	Kind          string  `json:"kind"`
	Generator     string  `json:"generator"`
	Length        int     `json:"length"`
	SHA256        string  `json:"sha256"`
	SeedNamespace string  `json:"seed_namespace"`
	SeedValue     uint64  `json:"seed_value"`
	BlockSize     int     `json:"block_size"`
	Blocks        []Block `json:"blocks"`
}

type GroundTruth struct {
	Schema string      `json:"schema"`
	Source SourceTruth `json:"source"`
}

// Build happens before any rendering. It returns both the exact source bytes and
// a structured truth object that later perception/rendering layers are judged against.
func Build(s spec.ExperimentSpec) (GroundTruth, labsource.Material, error) {
	material, err := labsource.Generate(s)
	if err != nil {
		return GroundTruth{}, labsource.Material{}, err
	}
	if len(material.Bytes) != s.Source.Length {
		return GroundTruth{}, labsource.Material{}, fmt.Errorf("ground truth: generated source length mismatch")
	}

	truth := GroundTruth{
		Schema: SchemaV1,
		Source: SourceTruth{
			Kind:          s.Source.Kind,
			Generator:     material.Generator,
			Length:        len(material.Bytes),
			SHA256:        hash(material.Bytes),
			SeedNamespace: material.SeedNamespace,
			SeedValue:     material.SeedValue,
			BlockSize:     DefaultBlockSize,
		},
	}
	for off, idx := 0, 0; off < len(material.Bytes); off, idx = off+DefaultBlockSize, idx+1 {
		end := off + DefaultBlockSize
		if end > len(material.Bytes) {
			end = len(material.Bytes)
		}
		truth.Source.Blocks = append(truth.Source.Blocks, Block{
			Index: idx, Offset: off, Length: end - off, SHA256: hash(material.Bytes[off:end]),
		})
	}
	return truth, material, nil
}

func Verify(t GroundTruth, b []byte) error {
	if t.Schema != SchemaV1 {
		return fmt.Errorf("ground truth: unsupported schema %q", t.Schema)
	}
	if len(b) != t.Source.Length {
		return fmt.Errorf("ground truth: length mismatch got=%d want=%d", len(b), t.Source.Length)
	}
	if got := hash(b); got != t.Source.SHA256 {
		return fmt.Errorf("ground truth: source sha mismatch got=%s want=%s", got, t.Source.SHA256)
	}
	for _, block := range t.Source.Blocks {
		if block.Offset < 0 || block.Length < 0 || block.Offset+block.Length > len(b) {
			return fmt.Errorf("ground truth: invalid block %d bounds", block.Index)
		}
		if got := hash(b[block.Offset : block.Offset+block.Length]); got != block.SHA256 {
			return fmt.Errorf("ground truth: block %d sha mismatch", block.Index)
		}
	}
	return nil
}

func hash(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
