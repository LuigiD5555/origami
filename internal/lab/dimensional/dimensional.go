package dimensional

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
)

type Scale string

const (
	ScaleMicro Scale = "MICRO"
	ScaleMeso  Scale = "MESO"
	ScaleMacro Scale = "MACRO"
)

type EvidenceStatus string

const (
	StatusUnknown EvidenceStatus = "UNKNOWN"
	StatusPass    EvidenceStatus = "PASS"
	StatusFail    EvidenceStatus = "FAIL"
)

type StateSpec struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	VisibleHint string `json:"visible_hint,omitempty"`
}

type DimensionSpec struct {
	ID               string          `json:"id"`
	Label            string          `json:"label"`
	Scale            Scale           `json:"scale"`
	FamilyIDs        []string        `json:"family_ids,omitempty"`
	Status           string          `json:"status"`
	Description      string          `json:"description,omitempty"`
	States           []StateSpec     `json:"states,omitempty"`
	Subdimensions    []DimensionSpec `json:"subdimensions,omitempty"`
	PerceptualStatus string          `json:"perceptual_status"`
}

type Registry struct {
	Schema     string          `json:"schema"`
	Version    string          `json:"version"`
	Dimensions []DimensionSpec `json:"dimensions"`
}

type Register struct {
	Schema      string            `json:"schema"`
	Registry    string            `json:"registry"`
	Assignments map[string]string `json:"assignments"`
}

type PairEvidence struct {
	A           string         `json:"a"`
	B           string         `json:"b"`
	Status      EvidenceStatus `json:"status"`
	EvidenceIDs []string       `json:"evidence_ids,omitempty"`
	Note        string         `json:"note,omitempty"`
}

type OrthogonalityMatrix struct {
	Schema string         `json:"schema"`
	Pairs  []PairEvidence `json:"pairs"`
}

type ScaleCapacity struct {
	Dimensions     []string `json:"dimensions"`
	AlphabetStates string   `json:"alphabet_states"`
	BitsPerUnit    float64  `json:"bits_per_unit"`
	Units          int      `json:"units"`
	TotalBits      float64  `json:"total_bits"`
}

type CapacityReport struct {
	RegistryVersion       string        `json:"registry_version"`
	Cells                 int           `json:"cells"`
	Groups                int           `json:"groups"`
	NominalMicro          ScaleCapacity `json:"nominal_micro"`
	NominalMeso           ScaleCapacity `json:"nominal_meso"`
	NominalMacro          ScaleCapacity `json:"nominal_macro"`
	NominalCarrierBits    float64       `json:"nominal_carrier_bits"`
	SafeMicro             ScaleCapacity `json:"safe_micro"`
	SafeMeso              ScaleCapacity `json:"safe_meso"`
	SafeMacro             ScaleCapacity `json:"safe_macro"`
	SafeCarrierBits       float64       `json:"safe_carrier_bits"`
	SafetyClaim           string        `json:"safety_claim"`
	PairwiseCandidateDims []string      `json:"pairwise_candidate_dimensions,omitempty"`
	Note                  string        `json:"note"`
}

func ValidateRegistry(r Registry) error {
	if r.Schema == "" || r.Version == "" {
		return errors.New("registry schema and version are required")
	}
	seen := map[string]bool{}
	var walk func([]DimensionSpec) error
	walk = func(ds []DimensionSpec) error {
		for _, d := range ds {
			if strings.TrimSpace(d.ID) == "" {
				return errors.New("dimension id is required")
			}
			if seen[d.ID] {
				return fmt.Errorf("duplicate dimension id %q", d.ID)
			}
			seen[d.ID] = true
			if d.Scale != ScaleMicro && d.Scale != ScaleMeso && d.Scale != ScaleMacro {
				return fmt.Errorf("dimension %s has invalid scale %q", d.ID, d.Scale)
			}
			hasStates := len(d.States) > 0
			hasChildren := len(d.Subdimensions) > 0
			if hasStates == hasChildren {
				return fmt.Errorf("dimension %s must contain either states or subdimensions, not both/neither", d.ID)
			}
			if hasStates {
				if len(d.States) < 2 {
					return fmt.Errorf("leaf dimension %s must have at least 2 states", d.ID)
				}
				stateSeen := map[string]bool{}
				for _, s := range d.States {
					if s.ID == "" || stateSeen[s.ID] {
						return fmt.Errorf("dimension %s has invalid or duplicate state id %q", d.ID, s.ID)
					}
					stateSeen[s.ID] = true
				}
			} else if err := walk(d.Subdimensions); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(r.Dimensions)
}

func LeafDimensions(r Registry) []DimensionSpec {
	out := []DimensionSpec{}
	var walk func([]DimensionSpec)
	walk = func(ds []DimensionSpec) {
		for _, d := range ds {
			if len(d.Subdimensions) == 0 {
				out = append(out, d)
				continue
			}
			walk(d.Subdimensions)
		}
	}
	walk(r.Dimensions)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func LeafMap(r Registry) map[string]DimensionSpec {
	m := map[string]DimensionSpec{}
	for _, d := range LeafDimensions(r) {
		m[d.ID] = d
	}
	return m
}

func ValidateRegister(r Registry, reg Register) error {
	if err := ValidateRegistry(r); err != nil {
		return err
	}
	leaves := LeafMap(r)
	for dimID, stateID := range reg.Assignments {
		d, ok := leaves[dimID]
		if !ok {
			return fmt.Errorf("unknown dimension %q", dimID)
		}
		found := false
		for _, s := range d.States {
			if s.ID == stateID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("dimension %s has unknown state %q", dimID, stateID)
		}
	}
	return nil
}

func pairKey(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "\x00" + b
}

func pairStatus(m OrthogonalityMatrix) map[string]EvidenceStatus {
	out := map[string]EvidenceStatus{}
	for _, p := range m.Pairs {
		if p.A == p.B || p.A == "" || p.B == "" {
			continue
		}
		out[pairKey(p.A, p.B)] = p.Status
	}
	return out
}

func MaxPairwiseCompatibleSet(r Registry, m OrthogonalityMatrix) []string {
	leaves := LeafDimensions(r)
	ids := make([]string, 0, len(leaves))
	for _, d := range leaves {
		if d.PerceptualStatus == "REJECTED" {
			continue
		}
		ids = append(ids, d.ID)
	}
	statuses := pairStatus(m)
	best := []string{}
	n := len(ids)
	if n > 22 {
		return best
	}
	for mask := 1; mask < (1 << n); mask++ {
		if bitsCount(mask) <= len(best) {
			continue
		}
		candidate := []string{}
		valid := true
		for i := 0; i < n && valid; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if mask&(1<<j) == 0 {
					continue
				}
				if statuses[pairKey(ids[i], ids[j])] != StatusPass {
					valid = false
					break
				}
			}
			if valid {
				candidate = append(candidate, ids[i])
			}
		}
		if valid && len(candidate) > len(best) {
			best = append([]string(nil), candidate...)
		}
	}
	sort.Strings(best)
	if len(best) < 2 {
		return nil
	}
	return best
}

func bitsCount(v int) int {
	n := 0
	for v != 0 {
		n += v & 1
		v >>= 1
	}
	return n
}

func productStates(r Registry, dimensionIDs []string) (*big.Int, error) {
	leaves := LeafMap(r)
	total := big.NewInt(1)
	for _, id := range dimensionIDs {
		d, ok := leaves[id]
		if !ok {
			return nil, fmt.Errorf("unknown dimension %q", id)
		}
		total.Mul(total, big.NewInt(int64(len(d.States))))
	}
	return total, nil
}

func bitsForBigInt(n *big.Int) float64 {
	if n.Sign() <= 0 || n.Cmp(big.NewInt(1)) == 0 {
		return 0
	}
	f, _ := new(big.Float).SetInt(n).Float64()
	return math.Log2(f)
}

func Capacity(r Registry, m OrthogonalityMatrix, cells, groups int) (CapacityReport, error) {
	if cells < 1 || groups < 1 {
		return CapacityReport{}, errors.New("cells and groups must be >= 1")
	}
	if err := ValidateRegistry(r); err != nil {
		return CapacityReport{}, err
	}
	leaves := LeafDimensions(r)
	promoted := []string{}
	for _, d := range leaves {
		if d.PerceptualStatus == "PROMOTED_SAFE" {
			promoted = append(promoted, d.ID)
		}
	}
	statuses := pairStatus(m)
	safe := []string{}
	for _, id := range promoted {
		okay := true
		for _, prior := range safe {
			if statuses[pairKey(id, prior)] != StatusPass {
				okay = false
				break
			}
		}
		if okay {
			safe = append(safe, id)
		}
	}

	byScale := func(ids []string, scale Scale, units int) (ScaleCapacity, error) {
		leafMap := LeafMap(r)
		selected := []string{}
		for _, id := range ids {
			if leafMap[id].Scale == scale {
				selected = append(selected, id)
			}
		}
		states := big.NewInt(1)
		var err error
		if len(selected) > 0 {
			states, err = productStates(r, selected)
			if err != nil {
				return ScaleCapacity{}, err
			}
		}
		bits := bitsForBigInt(states)
		return ScaleCapacity{Dimensions: selected, AlphabetStates: states.String(), BitsPerUnit: bits, Units: units, TotalBits: bits * float64(units)}, nil
	}

	all := make([]string, 0, len(leaves))
	for _, d := range leaves {
		all = append(all, d.ID)
	}
	nMicro, err := byScale(all, ScaleMicro, cells)
	if err != nil {
		return CapacityReport{}, err
	}
	nMeso, err := byScale(all, ScaleMeso, groups)
	if err != nil {
		return CapacityReport{}, err
	}
	nMacro, err := byScale(all, ScaleMacro, 1)
	if err != nil {
		return CapacityReport{}, err
	}
	sMicro, err := byScale(safe, ScaleMicro, cells)
	if err != nil {
		return CapacityReport{}, err
	}
	sMeso, err := byScale(safe, ScaleMeso, groups)
	if err != nil {
		return CapacityReport{}, err
	}
	sMacro, err := byScale(safe, ScaleMacro, 1)
	if err != nil {
		return CapacityReport{}, err
	}

	claim := "UNPROVEN"
	if len(safe) > 0 {
		claim = "PROMOTED_DIMENSIONS_ONLY; HIGHER_ORDER_COMPOUND_AND_RESTORATION_GATES_STILL_REQUIRED"
	}
	return CapacityReport{
		RegistryVersion:       r.Version,
		Cells:                 cells,
		Groups:                groups,
		NominalMicro:          nMicro,
		NominalMeso:           nMeso,
		NominalMacro:          nMacro,
		NominalCarrierBits:    nMicro.TotalBits + nMeso.TotalBits + nMacro.TotalBits,
		SafeMicro:             sMicro,
		SafeMeso:              sMeso,
		SafeMacro:             sMacro,
		SafeCarrierBits:       sMicro.TotalBits + sMeso.TotalBits + sMacro.TotalBits,
		SafetyClaim:           claim,
		PairwiseCandidateDims: MaxPairwiseCompatibleSet(r, m),
		Note:                  "Nominal capacity is a mathematical upper bound separated by visual scope: MICRO per cell, MESO per group, MACRO per carrier. SAFE_PSS requires Perception Lab evidence and does not inherit safety from primitive or pairwise tests alone.",
	}, nil
}

func MarshalCanonical(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func DefaultRegistry() Registry {
	state4 := func(ids ...string) []StateSpec {
		out := make([]StateSpec, 0, len(ids))
		for _, id := range ids {
			out = append(out, StateSpec{ID: id, Label: strings.ReplaceAll(id, "_", " ")})
		}
		return out
	}
	return Registry{
		Schema:  "ohf.dimensional-registry.v1",
		Version: "DVR-PSS-0.1.0",
		Dimensions: []DimensionSpec{
			{
				ID: "micro_intrinsic", Label: "Micro intrinsic register", Scale: ScaleMicro, Status: "REFERENCE_GROUP", PerceptualStatus: "GROUP_ONLY",
				Subdimensions: []DimensionSpec{
					{ID: "micro.outer_shape", Label: "Outer shape", Scale: ScaleMicro, FamilyIDs: []string{"F01"}, Status: "R310_PHYSICAL_REFERENCE", PerceptualStatus: "UNPROVEN_NATIVE", States: state4("circle", "square", "diamond", "triangle")},
					{ID: "micro.topology_holes", Label: "Topology holes", Scale: ScaleMicro, FamilyIDs: []string{"F02"}, Status: "R310_PHYSICAL_REFERENCE", PerceptualStatus: "UNPROVEN_NATIVE", States: state4("holes_0", "holes_1", "holes_2", "holes_3")},
					{ID: "micro.frame_count", Label: "Frame count", Scale: ScaleMicro, FamilyIDs: []string{"F04"}, Status: "R310_PHYSICAL_REFERENCE", PerceptualStatus: "UNPROVEN_NATIVE", States: state4("frames_0", "frames_1", "frames_2", "frames_3")},
				},
			},
			{
				ID: "micro_directional", Label: "Micro directional register", Scale: ScaleMicro, Status: "REFERENCE_GROUP", PerceptualStatus: "GROUP_ONLY",
				Subdimensions: []DimensionSpec{
					{ID: "micro.direction_beacon", Label: "Direction beacon", Scale: ScaleMicro, FamilyIDs: []string{"F03"}, Status: "R310_PHYSICAL_REFERENCE", PerceptualStatus: "UNPROVEN_NATIVE", States: state4("north", "east", "south", "west")},
				},
			},
			{
				ID: "meso_relational", Label: "Meso relational register", Scale: ScaleMeso, Status: "LAB_CANDIDATE_GROUP", PerceptualStatus: "GROUP_ONLY",
				Subdimensions: []DimensionSpec{
					{ID: "meso.ordering", Label: "Ordering", Scale: ScaleMeso, FamilyIDs: []string{"F28"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("row_major", "column_major", "clockwise", "radial")},
					{ID: "meso.grouping", Label: "Grouping", Scale: ScaleMeso, FamilyIDs: []string{"F29"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("singleton", "pair", "cluster", "ring")},
				},
			},
			{
				ID: "macro_gestalt", Label: "Macro gestalt register", Scale: ScaleMacro, Status: "LAB_CANDIDATE_GROUP", PerceptualStatus: "GROUP_ONLY",
				Subdimensions: []DimensionSpec{
					{ID: "macro.topology", Label: "Macro topology", Scale: ScaleMacro, FamilyIDs: []string{"F35", "F40"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("radial", "branching", "bilateral", "nested")},
					{ID: "macro.symmetry", Label: "Macro symmetry", Scale: ScaleMacro, FamilyIDs: []string{"F37"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("none", "bilateral", "rotational_4", "rotational_6")},
					{ID: "macro.orientation", Label: "Macro orientation", Scale: ScaleMacro, FamilyIDs: []string{"F03", "F39"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("north", "east", "south", "west")},
					{ID: "macro.region_structure", Label: "Macro region structure", Scale: ScaleMacro, FamilyIDs: []string{"F40", "F41"}, Status: "LAB_CANDIDATE", PerceptualStatus: "UNTESTED", States: state4("single_core", "dual_core", "ring_core", "distributed")},
				},
			},
		},
	}
}
