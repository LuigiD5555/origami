package baseline

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const ArtifactName = "ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip"

type Lock struct {
	Schema     string   `json:"schema"`
	Baseline   string   `json:"baseline"`
	Artifact   string   `json:"artifact"`
	SHA256     string   `json:"sha256"`
	Bytes      int64    `json:"bytes"`
	ZipEntries []string `json:"zip_entries"`
	TrustNote  string   `json:"trust_note"`
}

type Verification struct {
	Status       string
	ArtifactPath string
	SHA256       string
	Bytes        int64
	EntryCount   int
	LockPresent  bool
	Problems     []string
}

func artifactPath(root string) string {
	return filepath.Join(root, "baseline", "r310", ArtifactName)
}

func lockPath(root string) string {
	return filepath.Join(root, "baseline", "r310", "BASELINE_LOCK.json")
}

func Verify(root string) Verification {
	v := Verification{Status: "FAIL", ArtifactPath: artifactPath(root)}
	st, err := os.Stat(v.ArtifactPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			v.Problems = append(v.Problems, "MISSING_R310_SDK")
			return v
		}
		v.Problems = append(v.Problems, "STAT_ERROR: "+err.Error())
		return v
	}
	v.Bytes = st.Size()
	sha, err := fileSHA256(v.ArtifactPath)
	if err != nil {
		v.Problems = append(v.Problems, "HASH_ERROR: "+err.Error())
		return v
	}
	v.SHA256 = sha

	entries, err := zipEntries(v.ArtifactPath)
	if err != nil {
		v.Problems = append(v.Problems, "ZIP_INVALID: "+err.Error())
		return v
	}
	v.EntryCount = len(entries)

	requiredHints := []string{"cmd/", "internal/", "perception-lab/", "profiles/", "generated/"}
	for _, hint := range requiredHints {
		if !containsPrefix(entries, hint) {
			v.Problems = append(v.Problems, "MISSING_EXPECTED_TREE_HINT: "+hint)
		}
	}

	lockBytes, err := os.ReadFile(lockPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			v.Problems = append(v.Problems, "BASELINE_UNFROZEN")
			return v
		}
		v.Problems = append(v.Problems, "LOCK_READ_ERROR: "+err.Error())
		return v
	}
	v.LockPresent = true
	var l Lock
	if err := json.Unmarshal(lockBytes, &l); err != nil {
		v.Problems = append(v.Problems, "LOCK_INVALID: "+err.Error())
		return v
	}
	if l.Schema != "ohf.baseline-lock.v1" || l.Baseline != "R3.10" || l.Artifact != ArtifactName {
		v.Problems = append(v.Problems, "LOCK_IDENTITY_MISMATCH")
	}
	if l.SHA256 != v.SHA256 || l.Bytes != v.Bytes {
		v.Problems = append(v.Problems, "LOCK_HASH_MISMATCH")
	}
	if len(v.Problems) == 0 {
		v.Status = "PASS"
	}
	return v
}

func Freeze(root string) (Lock, error) {
	path := artifactPath(root)
	st, err := os.Stat(path)
	if err != nil {
		return Lock{}, err
	}
	sha, err := fileSHA256(path)
	if err != nil {
		return Lock{}, err
	}
	entries, err := zipEntries(path)
	if err != nil {
		return Lock{}, err
	}
	sort.Strings(entries)
	l := Lock{
		Schema:     "ohf.baseline-lock.v1",
		Baseline:   "R3.10",
		Artifact:   ArtifactName,
		SHA256:     sha,
		Bytes:      st.Size(),
		ZipEntries: entries,
		TrustNote:  "Operator confirmed this artifact is the intended historical R3.10 SDK. The continuity record does not provide an original SHA256, so authenticity cannot be cryptographically inferred from the record alone.",
	}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return Lock{}, err
	}
	b = append(b, '\n')
	if err := os.WriteFile(lockPath(root), b, 0o644); err != nil {
		return Lock{}, err
	}
	return l, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func zipEntries(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	out := make([]string, 0, len(r.File))
	for _, f := range r.File {
		name := strings.TrimPrefix(filepath.ToSlash(f.Name), "./")
		out = append(out, name)
	}
	return out, nil
}

func containsPrefix(entries []string, prefix string) bool {
	for _, e := range entries {
		if strings.HasPrefix(e, prefix) || strings.Contains(e, "/"+prefix) {
			return true
		}
	}
	return false
}

func (v Verification) Error() error {
	if v.Status == "PASS" {
		return nil
	}
	return fmt.Errorf("baseline verification failed: %s", strings.Join(v.Problems, ", "))
}
