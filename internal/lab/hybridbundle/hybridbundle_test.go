package hybridbundle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildCreatesOnlyDeclaredPublicInputs(t *testing.T) {
	root := t.TempDir()
	fixtures := map[string][]byte{
		"prompt.md": []byte("bootstrap"),
		"carrier.png": []byte("png-fixture"),
		"packet.json": []byte(`{"allowed_tools":["LOOKUP"]}`),
		"query.txt": []byte("question"),
	}
	for name, data := range fixtures {
		if err := os.WriteFile(filepath.Join(root, name), data, 0o644); err != nil { t.Fatal(err) }
	}
	bundle, err := Build(root, "run-1", filepath.Join(root, "prompt.md"), filepath.Join(root, "carrier.png"), filepath.Join(root, "packet.json"), filepath.Join(root, "query.txt"))
	if err != nil { t.Fatal(err) }
	if bundle.Schema != SchemaR0 || len(bundle.Files) != 4 { t.Fatalf("unexpected bundle: %+v", bundle) }
	public := filepath.Join(root, "runs", "run-1", filepath.FromSlash(bundle.PublicRelativePath))
	if err := VerifyPublic(public); err != nil { t.Fatal(err) }

	if err := os.WriteFile(filepath.Join(public, "oracle.json"), []byte(`{"answer":"hidden"}`), 0o644); err != nil { t.Fatal(err) }
	if err := VerifyPublic(public); err == nil {
		t.Fatal("expected undeclared private input to violate hybrid isolation")
	}
}
