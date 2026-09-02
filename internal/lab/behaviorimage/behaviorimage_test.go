package behaviorimage

import (
	"bytes"
	"encoding/json"
	"testing"
)

func fixture() Recipe {
	return Recipe{Schema: SchemaR0, ID: "find-number-r0", Boot: "read recipe then execute one step at a time", OutputContract: "NUMBER OR UNKNOWN", Steps: []Step{
		{ID: "locate", Capability: "LOCATE_REGION", Instruction: "find the region named TARGET", InputRef: "IMAGE", OutputRef: "REGION", OnSuccess: "read", OnUnknown: "UNKNOWN"},
		{ID: "read", Capability: "READ_SHORT_TEXT", Instruction: "read only the selected region", InputRef: "REGION", OutputRef: "TEXT", OnSuccess: "extract", OnUnknown: "UNKNOWN"},
		{ID: "extract", Capability: "EXTRACT_NUMBER", Instruction: "return the first relevant number", InputRef: "TEXT", OutputRef: "NUMBER", Halt: true, OnUnknown: "UNKNOWN"},
	}}
}

func TestRenderDeterministic(t *testing.T) {
	recipe := fixture()
	a, ma, err := Render(recipe, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	b, mb, err := Render(recipe, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) || ma.PNGSHA256 != mb.PNGSHA256 {
		t.Fatal("behavior image render must be deterministic")
	}
	if ma.StepCount != 3 || ma.RecipeSHA256 == "" || ma.PNGSHA256 == "" {
		t.Fatalf("bad manifest: %+v", ma)
	}
}

func TestValidateRejectsBackwardJump(t *testing.T) {
	r := fixture()
	r.Steps[1].OnSuccess = "locate"
	if err := Validate(r); err == nil {
		t.Fatal("backward jump should fail closed in R0")
	}
}
func TestValidateRejectsOversizedProgram(t *testing.T) {
	r := fixture()
	for len(r.Steps) <= DefaultMaxSteps {
		r.Steps = append(r.Steps, Step{ID: "x", Capability: "READ_SHORT_TEXT", Instruction: "read", Halt: true})
	}
	if err := Validate(r); err == nil {
		t.Fatal("oversized recipe must be rejected")
	}
}
func TestDecodeRecipeRoundtrip(t *testing.T) {
	r := fixture()
	raw, _ := json.Marshal(r)
	got, err := DecodeRecipe(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID || len(got.Steps) != len(r.Steps) {
		t.Fatalf("roundtrip mismatch: %+v", got)
	}
}
