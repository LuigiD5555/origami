package dimensional

import (
	"testing"
)

func TestDefaultRegistryAndNominalCapacity(t *testing.T) {
	r := DefaultRegistry()
	if err := ValidateRegistry(r); err != nil {
		t.Fatal(err)
	}
	leaves := LeafDimensions(r)
	if got, want := len(leaves), 10; got != want {
		t.Fatalf("leaf dimensions=%d want=%d", got, want)
	}
	report, err := Capacity(r, OrthogonalityMatrix{Schema: "ohf.orthogonality.v1"}, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if report.NominalMicro.AlphabetStates != "256" || report.NominalMicro.BitsPerUnit != 8 {
		t.Fatalf("micro nominal=%+v", report.NominalMicro)
	}
	if report.NominalMeso.AlphabetStates != "16" || report.NominalMeso.BitsPerUnit != 4 {
		t.Fatalf("meso nominal=%+v", report.NominalMeso)
	}
	if report.NominalMacro.AlphabetStates != "256" || report.NominalMacro.BitsPerUnit != 8 {
		t.Fatalf("macro nominal=%+v", report.NominalMacro)
	}
	if report.NominalCarrierBits != 20 {
		t.Fatalf("nominal carrier bits=%v", report.NominalCarrierBits)
	}
	if report.SafeCarrierBits != 0 || report.SafetyClaim != "UNPROVEN" {
		t.Fatalf("safe capacity must remain unproven: %+v", report)
	}
}

func TestValidateDVR(t *testing.T) {
	r := DefaultRegistry()
	reg := Register{
		Schema:   "ohf.dvr.v1",
		Registry: r.Version,
		Assignments: map[string]string{
			"micro.outer_shape": "triangle",
			"macro.orientation": "east",
			"meso.grouping":     "cluster",
			"micro.frame_count": "frames_2",
		},
	}
	if err := ValidateRegister(r, reg); err != nil {
		t.Fatal(err)
	}
	reg.Assignments["macro.orientation"] = "up-left"
	if err := ValidateRegister(r, reg); err == nil {
		t.Fatal("expected invalid state to fail")
	}
}

func TestOrthogonalityRequiresExplicitPass(t *testing.T) {
	r := DefaultRegistry()
	m := OrthogonalityMatrix{Schema: "ohf.orthogonality.v1", Pairs: []PairEvidence{
		{A: "macro.topology", B: "macro.symmetry", Status: StatusPass},
		{A: "macro.topology", B: "macro.orientation", Status: StatusPass},
		{A: "macro.symmetry", B: "macro.orientation", Status: StatusPass},
	}}
	got := MaxPairwiseCompatibleSet(r, m)
	if len(got) != 3 {
		t.Fatalf("candidate clique=%v", got)
	}
}

func TestPromotedSafeCapacityStillNeedsPairEvidence(t *testing.T) {
	r := DefaultRegistry()
	for i := range r.Dimensions {
		for j := range r.Dimensions[i].Subdimensions {
			d := &r.Dimensions[i].Subdimensions[j]
			if d.ID == "micro.outer_shape" || d.ID == "micro.direction_beacon" {
				d.PerceptualStatus = "PROMOTED_SAFE"
			}
		}
	}
	noPair, err := Capacity(r, OrthogonalityMatrix{Schema: "ohf.orthogonality.v1"}, 8, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(noPair.SafeMicro.Dimensions) != 1 {
		t.Fatalf("without pair evidence only one promoted micro dimension may be admitted, got %v", noPair.SafeMicro.Dimensions)
	}
	withPair, err := Capacity(r, OrthogonalityMatrix{Schema: "ohf.orthogonality.v1", Pairs: []PairEvidence{
		{A: "micro.outer_shape", B: "micro.direction_beacon", Status: StatusPass},
	}}, 8, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(withPair.SafeMicro.Dimensions) != 2 || withPair.SafeMicro.AlphabetStates != "16" || withPair.SafeMicro.BitsPerUnit != 4 {
		t.Fatalf("unexpected safe report: %+v", withPair)
	}
}
