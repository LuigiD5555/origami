package canonical

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/LuigiD5555/origami/internal/lab/spec"
)

func SpecBytes(s spec.ExperimentSpec) ([]byte, error) {
	return json.Marshal(s)
}

func SHA256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func SpecHash(s spec.ExperimentSpec) (string, error) {
	b, err := SpecBytes(s)
	if err != nil {
		return "", err
	}
	return SHA256Hex(b), nil
}
