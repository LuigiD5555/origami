package codec

import "fmt"

const (
	Profile3ID        = "origami.fixed-carrier.r2.profile-3"
	Profile3Width     = 640
	Profile3Height    = 640
	Profile3TargetPNG = 8192
	Profile3HardMax   = 512000
)

type Profile3Contract struct {
	ProfileID      string   `json:"profile_id"`
	Status         string   `json:"status"`
	Width          int      `json:"width"`
	Height         int      `json:"height"`
	TargetPNGBytes int      `json:"target_png_bytes"`
	HardMaxBytes   int      `json:"hard_max_bytes"`
	CoreCodecs     []string `json:"core_codecs"`
	NativeIndex    []string `json:"native_index_route"`
	WriteIndex     []string `json:"write_index_route"`
	RoundTripGate  string   `json:"roundtrip_gate"`
}

func ReferenceProfile3() Profile3Contract {
	return Profile3Contract{
		ProfileID:      Profile3ID,
		Status:         "EXPERIMENTAL_CONSTRUCTION_SPEC_ONLY",
		Width:          Profile3Width,
		Height:         Profile3Height,
		TargetPNGBytes: Profile3TargetPNG,
		HardMaxBytes:   Profile3HardMax,
		CoreCodecs:     []string{"S0", "E0", "S1", "E1", "S2", "E2"},
		NativeIndex:    []string{"T0", "T1", "S2", "T2"},
		WriteIndex:     []string{"SEMANTIC_IR", "E2", "T2_CONSTRUCTION_IR"},
		RoundTripGate:  "S2(E2(INDEX)) ~= INDEX",
	}
}

func ValidateProfile3(p Profile3Contract, r Registry) error {
	if p.ProfileID != Profile3ID {
		return fmt.Errorf("unexpected profile ID %q", p.ProfileID)
	}
	if p.Width != Profile3Width || p.Height != Profile3Height {
		return fmt.Errorf("profile-3 geometry drift: %dx%d", p.Width, p.Height)
	}
	if p.TargetPNGBytes != Profile3TargetPNG || p.HardMaxBytes != Profile3HardMax {
		return fmt.Errorf("profile-3 transport bounds drift")
	}
	if p.TargetPNGBytes > p.HardMaxBytes {
		return fmt.Errorf("profile-3 target exceeds hard maximum")
	}
	if err := r.Validate(); err != nil {
		return err
	}
	for _, id := range p.CoreCodecs {
		if _, ok := r.Lookup(id); !ok {
			return fmt.Errorf("profile-3 core codec %s is not declared", id)
		}
	}
	if p.RoundTripGate != "S2(E2(INDEX)) ~= INDEX" {
		return fmt.Errorf("profile-3 first roundtrip gate drift")
	}
	return nil
}
