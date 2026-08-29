package rendererfingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Compute(root, id, version string) (string, error) {
	var files []string
	switch id + "@" + version {
	case "lab-geometry@v1":
		files = []string{"internal/lab/geometry/geometry.go", "internal/lab/renderer/renderer.go", "internal/lab/renderer/labgeometry/labgeometry.go"}
	case "pdf-carrier@v1":
		files = []string{
			"internal/lab/renderer/renderer.go",
			"internal/lab/renderer/pdfcarrier/pdfcarrier.go",
			"internal/lab/canonicaldocument/document.go",
			"internal/lab/superindex/superindex.go",
			"internal/lab/capacity/capacity.go",
		}
	case "superindex-carrier@v1":
		files = []string{"internal/lab/renderer/renderer.go", "internal/lab/renderer/superindexcarrier/superindexcarrier.go", "internal/lab/canonicaldocument/document.go", "internal/lab/superindex/superindex.go", "internal/lab/capacity/capacity.go"}
	default:
		return "", nil
	}
	h := sha256.New()
	fmt.Fprintf(h, "renderer=%s@%s\ngo=%s\n", id, version, runtime.Version())
	for _, rel := range files {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return "", fmt.Errorf("renderer fingerprint %s: %w", rel, err)
		}
		fmt.Fprintf(h, "file=%s\nbytes=%d\n", rel, len(b))
		_, _ = h.Write(b)
		_, _ = h.Write([]byte{'\n'})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
