package promotion

import (
	"fmt"
	"strings"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

const (
	SchemaR1      = "origami.perception-promotion.r1"
	ObservationR1 = "origami.perception-observation.r1"
)

type EvidenceKind string

const (
	EvidenceMock      EvidenceKind = "MOCK"
	EvidenceRealModel EvidenceKind = "REAL_MODEL"
)

type T3Observation struct {
	CarrierID     string `json:"carrier_id"`
	StoreRoot     string `json:"store_root_sha256"`
	SourceSHA256  string `json:"source_sha256"`
	PageCount     uint32 `json:"page_count"`
	BlockCount    uint32 `json:"block_count"`
	DocumentCount uint32 `json:"document_count"`
	ObjectCount   uint32 `json:"object_count"`
}

type Observation struct {
	Schema       string        `json:"schema"`
	Model        string        `json:"model"`
	Trial        int           `json:"trial"`
	Transport    string        `json:"transport"`
	EvidenceKind EvidenceKind  `json:"evidence_kind"`
	BootText     []string      `json:"boot_text"`
	ProbeTop     string        `json:"probe_top"`
	ProbeBottom  string        `json:"probe_bottom"`
	ToolProtocol string        `json:"tool_protocol"`
	AddressABI   string        `json:"address_abi"`
	T3           *T3Observation `json:"t3,omitempty"`
}

type Check struct {
	ID       string `json:"id"`
	Pass     bool   `json:"pass"`
	Expected any    `json:"expected,omitempty"`
	Observed any    `json:"observed,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type Report struct {
	Schema                   string       `json:"schema"`
	Model                    string       `json:"model"`
	Trial                    int          `json:"trial"`
	Transport                string       `json:"transport"`
	EvidenceKind             EvidenceKind `json:"evidence_kind"`
	CarrierDigest            string       `json:"carrier_digest_sha256"`
	StoreRoot                string       `json:"store_root_sha256"`
	Checks                   []Check      `json:"checks"`
	HybridMechanicalPass     bool         `json:"hybrid_mechanical_pass"`
	NativeT3MechanicalPass   bool         `json:"native_t3_mechanical_pass"`
	HybridTrialPromotionOK   bool         `json:"hybrid_trial_promotion_eligible"`
	NativeT3TrialPromotionOK bool         `json:"native_t3_trial_promotion_eligible"`
	PromotionBoundary        string       `json:"promotion_boundary"`
}

func Evaluate(carrierPNG []byte, observation Observation) (Report, error) {
	if observation.Schema != ObservationR1 {
		return Report{}, fmt.Errorf("observation schema must be %q", ObservationR1)
	}
	if strings.TrimSpace(observation.Model) == "" || observation.Trial <= 0 {
		return Report{}, fmt.Errorf("model and positive trial are required")
	}
	if observation.EvidenceKind != EvidenceMock && observation.EvidenceKind != EvidenceRealModel {
		return Report{}, fmt.Errorf("unsupported evidence_kind %q", observation.EvidenceKind)
	}
	decoded, err := fixedcarrier.DecodePNG(carrierPNG)
	if err != nil {
		return Report{}, fmt.Errorf("decode canonical carrier: %w", err)
	}
	probeTop, probeBottom, err := fixedcarrier.VisualProbeFromPNG(carrierPNG)
	if err != nil {
		return Report{}, fmt.Errorf("read canonical visual probes: %w", err)
	}

	report := Report{
		Schema:            SchemaR1 + ".report",
		Model:             observation.Model,
		Trial:             observation.Trial,
		Transport:         observation.Transport,
		EvidenceKind:      observation.EvidenceKind,
		CarrierDigest:     decoded.CarrierDigest,
		StoreRoot:         decoded.StoreRoot,
		PromotionBoundary: "SINGLE_TRIAL_ONLY_NOT_STACK_SUPPORT",
	}

	bootPass := equalNormalizedLines(observation.BootText, decoded.BootText)
	report.Checks = append(report.Checks, Check{ID: "T0_BOOT", Pass: bootPass, Expected: decoded.BootText, Observed: observation.BootText, Reason: reason(bootPass, "boot text mismatch")})

	topPass := strings.TrimSpace(observation.ProbeTop) == probeTop && probeTop == decoded.VisualProbe
	bottomPass := strings.TrimSpace(observation.ProbeBottom) == probeBottom && probeBottom == decoded.VisualProbe
	agreePass := strings.TrimSpace(observation.ProbeTop) != "" && strings.TrimSpace(observation.ProbeTop) == strings.TrimSpace(observation.ProbeBottom)
	report.Checks = append(report.Checks,
		Check{ID: "T1_PROBE_TOP", Pass: topPass, Expected: probeTop, Observed: observation.ProbeTop, Reason: reason(topPass, "top visual probe mismatch")},
		Check{ID: "T1_PROBE_BOTTOM", Pass: bottomPass, Expected: probeBottom, Observed: observation.ProbeBottom, Reason: reason(bottomPass, "bottom visual probe mismatch")},
		Check{ID: "T1_PROBE_DUPLICATION", Pass: agreePass, Expected: probeTop, Observed: observation.ProbeTop + "/" + observation.ProbeBottom, Reason: reason(agreePass, "duplicated visual probes do not agree")},
	)

	toolPass := strings.TrimSpace(observation.ToolProtocol) == decoded.ToolProtocol
	addressPass := strings.TrimSpace(observation.AddressABI) == decoded.AddressABI
	report.Checks = append(report.Checks,
		Check{ID: "TOOL_PROTOCOL", Pass: toolPass, Expected: decoded.ToolProtocol, Observed: observation.ToolProtocol, Reason: reason(toolPass, "tool protocol mismatch")},
		Check{ID: "ADDRESS_ABI", Pass: addressPass, Expected: decoded.AddressABI, Observed: observation.AddressABI, Reason: reason(addressPass, "address ABI mismatch")},
	)

	report.HybridMechanicalPass = bootPass && topPass && bottomPass && agreePass && toolPass && addressPass

	t3Pass := false
	if observation.T3 != nil {
		t3Pass = observation.T3.CarrierID == decoded.CarrierID &&
			observation.T3.StoreRoot == decoded.StoreRoot &&
			observation.T3.SourceSHA256 == decoded.SourceSHA256 &&
			observation.T3.PageCount == decoded.PageCount &&
			observation.T3.BlockCount == decoded.BlockCount &&
			observation.T3.DocumentCount == decoded.DocumentCount &&
			observation.T3.ObjectCount == decoded.ObjectCount
		report.Checks = append(report.Checks, Check{ID: "T3_MACHINE_RECORD", Pass: t3Pass, Expected: T3Observation{CarrierID: decoded.CarrierID, StoreRoot: decoded.StoreRoot, SourceSHA256: decoded.SourceSHA256, PageCount: decoded.PageCount, BlockCount: decoded.BlockCount, DocumentCount: decoded.DocumentCount, ObjectCount: decoded.ObjectCount}, Observed: *observation.T3, Reason: reason(t3Pass, "T3 machine record mismatch")})
	} else {
		report.Checks = append(report.Checks, Check{ID: "T3_MACHINE_RECORD", Pass: false, Reason: "not observed; Native T3 remains a separate gate"})
	}
	report.NativeT3MechanicalPass = report.HybridMechanicalPass && t3Pass
	real := observation.EvidenceKind == EvidenceRealModel
	report.HybridTrialPromotionOK = real && report.HybridMechanicalPass
	report.NativeT3TrialPromotionOK = real && report.NativeT3MechanicalPass
	return report, nil
}

func equalNormalizedLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if normalizeLine(a[i]) != normalizeLine(b[i]) {
			return false
		}
	}
	return true
}

func normalizeLine(value string) string {
	return strings.Join(strings.Fields(strings.ToUpper(strings.TrimSpace(value))), " ")
}

func reason(pass bool, failure string) string {
	if pass {
		return ""
	}
	return failure
}
