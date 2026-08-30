package selfcontainedbook

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

func TestPackRenderDecodeRoundTrip(t *testing.T) {
	source := bytes.Repeat([]byte("Origami semantic source.\n"), 200)
	compressed := []byte("synthetic-compressed-residual-for-transport-test")
	prompt := sha256.Sum256([]byte("master prompt"))
	promptHex := hex.EncodeToString(prompt[:])
	payload, header, err := Pack(source, compressed, 12, 34, "test", promptHex)
	if err != nil { t.Fatal(err) }
	h := sha256.Sum256(source)
	hash := hex.EncodeToString(h[:])
	meta := fixedcarrier.Metadata{CarrierID: hash, StoreRoot: hash, SourceSHA256: hash, PageCount: 12, BlockCount: 34, DocumentCount: 1, ObjectCount: 47, GraphSignature: []byte("self-contained-book-r0-test")}
	controlPNG, _, err := fixedcarrier.Render(meta); if err != nil { t.Fatal(err) }
	rendered, err := Render(controlPNG, payload, 640); if err != nil { t.Fatal(err) }
	decodedPayload, _, err := DecodePNG(rendered.PNG); if err != nil { t.Fatal(err) }
	decodedHeader, gotCompressed, err := Unpack(decodedPayload); if err != nil { t.Fatal(err) }
	if decodedHeader.SourceSHA256 != header.SourceSHA256 { t.Fatalf("source hash mismatch: %s != %s", decodedHeader.SourceSHA256, header.SourceSHA256) }
	if decodedHeader.MasterPromptSHA256 != promptHex { t.Fatalf("master prompt hash mismatch: %s != %s", decodedHeader.MasterPromptSHA256, promptHex) }
	if !bytes.Equal(gotCompressed, compressed) { t.Fatal("compressed residual mismatch") }
}

func TestUnpackRejectsCorruption(t *testing.T) {
	source := []byte("source")
	compressed := []byte("compressed")
	prompt := sha256.Sum256([]byte("master prompt"))
	payload, _, err := Pack(source, compressed, 1, 1, "test", hex.EncodeToString(prompt[:])); if err != nil { t.Fatal(err) }
	payload[len(payload)-1] ^= 0xff
	if _, _, err := Unpack(payload); err == nil { t.Fatal("expected residual hash mismatch") }
}
