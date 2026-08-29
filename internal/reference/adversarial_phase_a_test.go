package reference

import "testing"

// Phase A intentionally records unsupported semantic capabilities before their
// implementation. These tests must remain green while reporting UNSUPPORTED;
// they become real semantic assertions only when the corresponding capability
// is restored. This prevents fake PASSes and preserves the R0 oracle baseline.

type adversarialCapability struct {
	ID        string
	Supported bool
	Reason    string
}

func requireUnsupported(t *testing.T, c adversarialCapability) {
	t.Helper()
	if c.Supported {
		t.Fatalf("%s unexpectedly marked supported before its semantic gate exists", c.ID)
	}
	if c.Reason == "" {
		t.Fatalf("%s must record why the current engine is unsupported", c.ID)
	}
	t.Logf("%s: UNSUPPORTED: %s", c.ID, c.Reason)
}

func TestADV001HigherOrderInteractionUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{
		ID: "ADV-001",
		Reason: "reference engine represents only pairwise Relation entries and has no participating-set/higher-order condition model",
	})
}

func TestADV002ContextEvolutionUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{
		ID: "ADV-002",
		Reason: "Experiment and Execute have no first-class Context participating in F(S_t,C_t,R)",
	})
}

func TestADV003RulesAsDataUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{
		ID: "ADV-003",
		Reason: "transition semantics are hidden in Go step/contradiction control flow rather than finite declarative rules with stable IDs and provenance",
	})
}

func TestADV004SemanticStatusesUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{
		ID: "ADV-004",
		Reason: "State is map[string]string and does not model ABSENT, UNKNOWN, INHIBITED and CANCELLED as distinct semantic statuses with evidence",
	})
}

func TestADV005ObservationIsNotTransitionUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{
		ID: "ADV-005",
		Reason: "reference engine emits transition traces but has no Observation Contract runtime over immutable pre-observation evidence",
	})
}
