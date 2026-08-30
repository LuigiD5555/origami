package writer

import "github.com/LuigiD5555/origami/internal/core"

const (
	SchemaR0          = "origami.writer.r0"
	RequestSchemaR0   = SchemaR0 + ".request"
	PlanSchemaR0      = SchemaR0 + ".construction-plan"
	CanonicalProfileR0 = "origami.canonical-aesthetic.r0"
)

type Intent string

const (
	IntentIdentity      Intent = "IDENTITY"
	IntentHierarchy     Intent = "HIERARCHY"
	IntentRelation      Intent = "RELATION"
	IntentDependency    Intent = "DEPENDENCY"
	IntentScope         Intent = "SCOPE"
	IntentImportance    Intent = "IMPORTANCE"
	IntentState         Intent = "STATE"
	IntentUncertainty   Intent = "UNCERTAINTY"
	IntentEvidence      Intent = "EVIDENCE"
	IntentVerification  Intent = "VERIFICATION"
	IntentAddressability Intent = "ADDRESSABILITY"
	IntentTemporalOrder Intent = "TEMPORAL_ORDER"
	IntentRedundancy    Intent = "REDUNDANCY"
	IntentRouting       Intent = "ROUTING"
)

type Request struct {
	Schema         string             `json:"schema"`
	SourceIdentity string             `json:"source_identity"`
	SemanticGraph core.SemanticGraph `json:"semantic_graph"`
	IntentHints    []IntentHint       `json:"intent_hints,omitempty"`
}

type IntentHint struct {
	Address string `json:"address"`
	Intent  Intent `json:"intent"`
	Note    string `json:"note,omitempty"`
}

type VisualIntent struct {
	Address       string `json:"address"`
	Intent        Intent `json:"intent"`
	Dimension     string `json:"dimension"`
	Primitive     string `json:"primitive"`
	SemanticRole  string `json:"semantic_role"`
	Authority     string `json:"authority"`
}

type RosettaEntry struct {
	VisualToken  string `json:"visual_token"`
	SemanticRole string `json:"semantic_role"`
	Dimension    string `json:"dimension"`
}

type Rosetta struct {
	ProfileID        string         `json:"profile_id"`
	ProfileVersion   string         `json:"profile_version"`
	ActiveDimensions []string       `json:"active_dimensions"`
	Entries          []RosettaEntry `json:"entries"`
}

type Section struct {
	ID       string `json:"id"`
	Required bool   `json:"required"`
	Purpose  string `json:"purpose"`
}

type VerificationPlan struct {
	SemanticRoundtrip bool     `json:"semantic_roundtrip"`
	ExactFields       []string `json:"exact_fields"`
	RequiredChecks    []string `json:"required_checks"`
	FailureState      string   `json:"failure_state"`
}

type ConstructionPlan struct {
	Schema             string             `json:"schema"`
	ProfileID          string             `json:"profile_id"`
	SourceIdentity     string             `json:"source_identity"`
	SemanticCommitment string             `json:"semantic_commitment"`
	VisualIntents      []VisualIntent     `json:"visual_intents"`
	Rosetta            Rosetta            `json:"rosetta"`
	Sections           []Section          `json:"sections"`
	Verification       VerificationPlan   `json:"verification_plan"`
	CanonicalGraph     core.SemanticGraph `json:"canonical_graph"`
	Invariants         []string           `json:"invariants"`
}
