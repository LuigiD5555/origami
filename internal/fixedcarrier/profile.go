package fixedcarrier

// Profile describes the frozen physical and semantic layout of a Fixed Carrier.
// It is intentionally independent from the size of the corpus behind Origami.
type Profile struct {
	ID             string       `json:"id"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`
	FixedPNGBytes  int          `json:"fixed_png_bytes"`
	MaxPNGBytes    int          `json:"max_png_bytes"`
	BootText       []string     `json:"boot_text"`
	Zones          []Zone       `json:"zones"`
	FamilyBindings []FamilyBind `json:"family_bindings"`
}

type Zone struct {
	ID      string `json:"id"`
	Purpose string `json:"purpose"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	W       int    `json:"w"`
	H       int    `json:"h"`
}

type FamilyBind struct {
	Family string `json:"family"`
	Role   string `json:"role"`
	Zone   string `json:"zone"`
	Status string `json:"status"`
}

var profileR2 = Profile{
	ID:            ProfileID,
	Width:         640,
	Height:        640,
	FixedPNGBytes: 8192,
	MaxPNGBytes:   512000,
	BootText:      append([]string(nil), bootText...),
	Zones: []Zone{
		{ID: "T0", Purpose: "plaintext semantic-first bootloader", X: 18, Y: 16, W: 604, H: 112},
		{ID: "T1", Purpose: "Rosetta and duplicated visual-presence challenge", X: 24, Y: 132, W: 592, H: 92},
		{ID: "T2", Purpose: "VLM-readable semantic superindex and navigation entry points", X: 18, Y: 214, W: 604, H: 82},
		{ID: "T3", Purpose: "dense deterministic control/exact record plus state map", X: 18, Y: 296, W: 604, H: 252},
		{ID: "VERIFY", Purpose: "duplicated challenge and verification footer", X: 24, Y: 552, W: 592, H: 72},
	},
	FamilyBindings: []FamilyBind{
		{Family: "F01", Role: "semantic entry kind", Zone: "T2/T3", Status: "navigation_core_candidate"},
		{Family: "F05", Role: "binary/semantic state redundancy", Zone: "T1/T3", Status: "candidate_auxiliary"},
		{Family: "F09", Role: "macro/meso/micro scale class", Zone: "T2", Status: "candidate_auxiliary"},
		{Family: "F21", Role: "semantic adjacency", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F22", Role: "containment", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F24", Role: "graph connection", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F25", Role: "edge direction", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F29", Role: "cluster grouping", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F32", Role: "hierarchy depth", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F34", Role: "graph degree", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F39", Role: "spatial zone", Zone: "T0/T1/T2/T3", Status: "navigation_core_candidate"},
		{Family: "F40", Role: "macro semantic region", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F41", Role: "resolution transition", Zone: "T2", Status: "navigation_core_candidate"},
		{Family: "F42", Role: "verification status redundancy", Zone: "VERIFY", Status: "navigation_core_candidate"},
	},
}

func R2Profile() Profile { return profileR2 }
