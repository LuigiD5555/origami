package receiver

import "testing"

func baseSpec(bindings []SymbolBinding) Spec {
	return Spec{
		ContractID: ContractID,
		Boot: BootSpec{
			Version:          "r0",
			EntrySymbol:      "BOOT",
			StartState:       "s0",
			RosettaRef:       "ROSETTA",
			ProgramRef:       "PROGRAM",
			IndexRef:         "INDEX",
			MemoryRef:        "MEMORY",
			VerificationRef:  "VERIFICATION",
			RequiredSections: []string{"BOOT", "ROSETTA", "PROGRAM", "INDEX", "MEMORY", "VERIFICATION"},
			WorkingWindow:    4000,
		},
		Rosetta: bindings,
		Program: []MicroRule{
			{ID: "m0", State: "s0", Token: "OPEN", NextState: "s1"},
			{ID: "m1", State: "s1", Token: "VALUE", Emit: "AMBER-10593", NextState: "done"},
			{ID: "m2", State: "done", Token: "HALT", NextState: "done", Halt: true},
		},
	}
}

func TestRosettaPermutationPreservesSemantics(t *testing.T) {
	a := baseSpec([]SymbolBinding{{Physical: "triangle", Semantic: "OPEN"}, {Physical: "circle", Semantic: "VALUE"}, {Physical: "square", Semantic: "HALT"}})
	b := baseSpec([]SymbolBinding{{Physical: "circle", Semantic: "OPEN"}, {Physical: "square", Semantic: "VALUE"}, {Physical: "triangle", Semantic: "HALT"}})

	traceA, err := Execute(a, []string{"triangle", "circle", "square"})
	if err != nil { t.Fatal(err) }
	traceB, err := Execute(b, []string{"circle", "square", "triangle"})
	if err != nil { t.Fatal(err) }

	if traceA.Unknown || traceB.Unknown || !traceA.Halted || !traceB.Halted {
		t.Fatalf("unexpected traces: A=%+v B=%+v", traceA, traceB)
	}
	if len(traceA.Outputs) != 1 || len(traceB.Outputs) != 1 || traceA.Outputs[0] != traceB.Outputs[0] {
		t.Fatalf("carrier-local symbol permutations changed semantics: A=%v B=%v", traceA.Outputs, traceB.Outputs)
	}
}

func TestUnknownPhysicalSymbolFailsClosed(t *testing.T) {
	s := baseSpec([]SymbolBinding{{Physical: "triangle", Semantic: "OPEN"}, {Physical: "circle", Semantic: "VALUE"}, {Physical: "square", Semantic: "HALT"}})
	trace, err := Execute(s, []string{"triangle", "unknown-symbol"})
	if err != nil { t.Fatal(err) }
	if !trace.Unknown {
		t.Fatalf("expected UNKNOWN, got %+v", trace)
	}
	if len(trace.Outputs) != 0 {
		t.Fatalf("unknown path must not invent output: %v", trace.Outputs)
	}
}

func TestValidateRejectsHardCodedAmbiguity(t *testing.T) {
	s := baseSpec([]SymbolBinding{{Physical: "triangle", Semantic: "OPEN"}, {Physical: "triangle", Semantic: "VALUE"}, {Physical: "square", Semantic: "HALT"}})
	if err := Validate(s); err == nil {
		t.Fatal("expected duplicate physical symbol to be rejected")
	}
}
