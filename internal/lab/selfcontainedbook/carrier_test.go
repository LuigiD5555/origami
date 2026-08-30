package selfcontainedbook

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

func TestPackRenderDecodeRoundTripIncludesMasterPrompt(t *testing.T) {
	source := bytes.Repeat([]byte("Origami semantic source.\n"), 200)
	masterPrompt := []byte("# Origami Master Prompt\nREAD BOOT ROSETTA INDEX MEMORY VERIFY\n")
	compressed := []byte("synthetic-compressed-residual-for-transport-test")
	payload, header, err := Pack(source, masterPrompt, compressed, 12, 34, "test")
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(source)
	hash := hex.EncodeToString(h[:])
	meta := fixedcarrier.Metadata{
		CarrierID: hash, StoreRoot: hash, SourceSHA256: hash,
		PageCount: 12, BlockCount: 34, DocumentCount: 1, ObjectCount: 47,
		GraphSignature: []byte("self-contained-book-r0-test"),
	}
	controlPNG, _, err := fixedcarrier.Render(meta)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(controlPNG, payload, DefaultDataWidth)
	if err != nil {
		t.Fatal(err)
	}
	if !rendered.VisibleMasterBootstrap {
		t.Fatal("expected visible Master Prompt bootstrap panel")
	}
	decodedPayload, _, err := DecodePNG(rendered.PNG)
	if err != nil {
		t.Fatal(err)
	}
	decodedHeader, gotPrompt, gotCompressed, err := Unpack(decodedPayload)
	if err != nil {
		t.Fatal(err)
	}
	if decodedHeader.SourceSHA256 != header.SourceSHA256 {
		t.Fatalf("source hash mismatch: %s != %s", decodedHeader.SourceSHA256, header.SourceSHA256)
	}
	if decodedHeader.MasterPromptSHA256 != header.MasterPromptSHA256 {
		t.Fatalf("Master Prompt hash mismatch: %s != %s", decodedHeader.MasterPromptSHA256, header.MasterPromptSHA256)
	}
	if decodedHeader.MasterPromptBytes != uint32(len(masterPrompt)) {
		t.Fatalf("Master Prompt size mismatch: %d != %d", decodedHeader.MasterPromptBytes, len(masterPrompt))
	}
	if !bytes.Equal(gotPrompt, masterPrompt) {
		t.Fatal("Master Prompt bytes mismatch")
	}
	if !bytes.Equal(gotCompressed, compressed) {
		t.Fatal("compressed residual mismatch")
	}
}

func TestUnpackRejectsPromptCorruption(t *testing.T) {
	source := []byte("source")
	masterPrompt := []byte("master prompt exact bytes")
	compressed := []byte("compressed")
	payload, _, err := Pack(source, masterPrompt, compressed, 1, 1, "test")
	if err != nil {
		t.Fatal(err)
	}
	payload[HeaderBytes+3] ^= 0xff
	if _, _, _, err := Unpack(payload); err == nil {
		t.Fatal("expected Master Prompt hash mismatch")
	}
}

func TestUnpackRejectsResidualCorruption(t *testing.T) {
	source := []byte("source")
	masterPrompt := []byte("master prompt exact bytes")
	compressed := []byte("compressed")
	payload, _, err := Pack(source, masterPrompt, compressed, 1, 1, "test")
	if err != nil {
		t.Fatal(err)
	}
	payload[len(payload)-1] ^= 0xff
	if _, _, _, err := Unpack(payload); err == nil {
		t.Fatal("expected residual hash mismatch")
	}
}
