package promotion

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

func fixture(t *testing.T) ([]byte, fixedcarrier.Decoded, string, string) {
	t.Helper()
	meta := fixedcarrier.Metadata{
		CarrierID: strings.Repeat("11", 32),
		StoreRoot: strings.Repeat("22", 32),
		SourceSHA256: strings.Repeat("33", 32),
		PageCount: 124,
		BlockCount: 1255,
		DocumentCount: 5,
		ObjectCount: 1509,
		GraphSignature: []byte(`{"s":"origami.virtual-memory.r0","n":1509}`),
	}
	png, decoded, err := fixedcarrier.Render(meta)
	if err != nil { t.Fatal(err) }
	top, bottom, err := fixedcarrier.VisualProbeFromPNG(png)
	if err != nil { t.Fatal(err) }
	return png, decoded, top, bottom
}

func fullObservation(decoded fixedcarrier.Decoded, top, bottom string, kind EvidenceKind) Observation {
	return Observation{
		Schema: ObservationR1,
		Model: "model-a",
		Trial: 1,
		Transport: "original",
		EvidenceKind: kind,
		BootText: append([]string(nil), decoded.BootText...),
		ProbeTop: top,
		ProbeBottom: bottom,
		ToolProtocol: decoded.ToolProtocol,
		AddressABI: decoded.AddressABI,
		T3: &T3Observation{CarrierID: decoded.CarrierID, StoreRoot: decoded.StoreRoot, SourceSHA256: decoded.SourceSHA256, PageCount: decoded.PageCount, BlockCount: decoded.BlockCount, DocumentCount: decoded.DocumentCount, ObjectCount: decoded.ObjectCount},
	}
}

func TestMockCanPassMechanicallyButNeverPromote(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	report, err := Evaluate(png, fullObservation(decoded, top, bottom, EvidenceMock))
	if err != nil { t.Fatal(err) }
	if !report.HybridMechanicalPass || !report.NativeT3MechanicalPass { t.Fatalf("mock should exercise full mechanical path: %+v", report) }
	if report.HybridTrialPromotionOK || report.NativeT3TrialPromotionOK { t.Fatalf("mock evidence must never promote: %+v", report) }
}

func TestRealModelHybridCanPassWithoutNativeT3(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	obs := fullObservation(decoded, top, bottom, EvidenceRealModel)
	obs.T3 = nil
	report, err := Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if !report.HybridMechanicalPass || !report.HybridTrialPromotionOK { t.Fatalf("hybrid gate should pass independently of T3: %+v", report) }
	if report.NativeT3MechanicalPass || report.NativeT3TrialPromotionOK { t.Fatalf("Native T3 must remain independent: %+v", report) }
}

func TestRealModelNativeT3TrialRequiresAllT3Fields(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	obs := fullObservation(decoded, top, bottom, EvidenceRealModel)
	report, err := Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if !report.NativeT3TrialPromotionOK { t.Fatalf("complete real-model T3 observation should be trial-eligible: %+v", report) }
	obs.T3.StoreRoot = strings.Repeat("44", 32)
	report, err = Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if report.NativeT3MechanicalPass || report.NativeT3TrialPromotionOK { t.Fatalf("wrong T3 root must fail Native gate: %+v", report) }
	if !report.HybridMechanicalPass { t.Fatal("T3 failure must not erase Hybrid mechanical success") }
}

func TestWrongProbeProvesImageChallengeFailed(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	obs := fullObservation(decoded, top, bottom, EvidenceRealModel)
	obs.ProbeBottom = "00000000"
	report, err := Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if report.HybridMechanicalPass || report.HybridTrialPromotionOK { t.Fatalf("wrong duplicated visual probe must fail Hybrid trial: %+v", report) }
}

func TestBootWhitespaceIsNormalizedButContentIsNotParaphrased(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	obs := fullObservation(decoded, top, bottom, EvidenceRealModel)
	obs.BootText[0] = "  origami   fixed carrier r2  "
	report, err := Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if !report.HybridMechanicalPass { t.Fatalf("benign whitespace/case should normalize: %+v", report) }
	obs.BootText[0] = "some visual memory image"
	report, err = Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if report.HybridMechanicalPass { t.Fatal("semantic paraphrase must not pass deterministic BOOT check") }
}

func TestAddressAndToolProtocolsAreExactGates(t *testing.T) {
	png, decoded, top, bottom := fixture(t)
	obs := fullObservation(decoded, top, bottom, EvidenceRealModel)
	obs.ToolProtocol = "other.tools.r1"
	report, err := Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if report.HybridMechanicalPass { t.Fatal("wrong tool protocol passed") }
	obs = fullObservation(decoded, top, bottom, EvidenceRealModel)
	obs.AddressABI = "other-address.r1"
	report, err = Evaluate(png, obs)
	if err != nil { t.Fatal(err) }
	if report.HybridMechanicalPass { t.Fatal("wrong address ABI passed") }
}
