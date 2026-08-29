package hybridruntime

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/hybridcarrier"
	"github.com/LuigiD5555/origami/internal/receiver"
)

func testReceiver() receiver.Spec {
	return receiver.Spec{
		ContractID: receiver.ContractID,
		Boot: receiver.BootSpec{
			Version: "r0", EntrySymbol: "BOOT", StartState: "boot",
			RosettaRef: "ROSETTA", ProgramRef: "PROGRAM", IndexRef: "INDEX", MemoryRef: "MEMORY", VerificationRef: "VERIFICATION",
			RequiredSections: []string{"BOOT", "ROSETTA", "PROGRAM", "INDEX", "MEMORY", "VERIFICATION"},
			WorkingWindow: 4000,
		},
		Rosetta: []receiver.SymbolBinding{
			{Physical: "triangle", Semantic: "OPEN"},
			{Physical: "circle", Semantic: "VALUE"},
		},
		Program: []receiver.MicroRule{
			{ID: "m0", State: "boot", Token: "OPEN", NextState: "memory"},
			{ID: "m1", State: "memory", Token: "VALUE", Emit: "READY", NextState: "done", Halt: true},
		},
	}
}

func TestRuntimeLookupFollowAndVerifyFromPNG(t *testing.T) {
	envelope, err := hybridcarrier.BuildEnvelope(
		testReceiver(),
		[]hybridcarrier.IndexEntry{{Key: "K7F91", Address: "mem:1"}},
		[]hybridcarrier.MemoryEntry{
			{Address: "mem:1", Value: "VIOLET-48271", Links: map[string][]string{"depends": []string{"mem:2"}}},
			{Address: "mem:2", Value: "DELTA-77102", Links: map[string][]string{"depends": []string{"mem:3"}}},
			{Address: "mem:3", Value: "AMBER-10593"},
		},
	)
	if err != nil { t.Fatal(err) }
	rendered, err := hybridcarrier.Render(envelope, 0, 0, 0, 0)
	if err != nil { t.Fatal(err) }

	runtime, err := OpenPNG(rendered.PNG, Transport{})
	if err != nil { t.Fatal(err) }
	metrics := runtime.OpenMetrics()
	if metrics.MemoryEntries != 3 || metrics.IndexEntries != 1 || metrics.BuildTouched != 4 {
		t.Fatalf("unexpected open metrics: %+v", metrics)
	}

	boot := runtime.Boot()
	if boot.Boot == nil || boot.Boot.RosettaRef != "ROSETTA" { t.Fatalf("unexpected boot: %+v", boot) }

	lookup := runtime.Lookup("K7F91")
	if lookup.Unknown || len(lookup.Entries) != 1 || lookup.Entries[0].Value != "VIOLET-48271" {
		t.Fatalf("unexpected lookup: %+v", lookup)
	}

	follow := runtime.Follow("K7F91", "depends", 2)
	if follow.Unknown || len(follow.Entries) != 3 || follow.Entries[2].Value != "AMBER-10593" {
		t.Fatalf("unexpected follow: %+v", follow)
	}
	if follow.Metrics.Touched >= metrics.MemoryEntries+metrics.IndexEntries+10 {
		t.Fatalf("unexpected traversal cost: %+v", follow.Metrics)
	}

	verify := runtime.Verify()
	if verify.Verified == nil || !*verify.Verified { t.Fatalf("unexpected verify: %+v", verify) }
}

func TestRuntimeMissingLookupReturnsUnknownWithoutScan(t *testing.T) {
	envelope, err := hybridcarrier.BuildEnvelope(
		testReceiver(), nil,
		[]hybridcarrier.MemoryEntry{{Address: "mem:1", Value: "A"}, {Address: "mem:2", Value: "B"}},
	)
	if err != nil { t.Fatal(err) }
	rendered, err := hybridcarrier.Render(envelope, 0, 0, 0, 0)
	if err != nil { t.Fatal(err) }
	runtime, err := OpenPNG(rendered.PNG, Transport{})
	if err != nil { t.Fatal(err) }
	result := runtime.Lookup("missing")
	if !result.Unknown || result.Metrics.Touched != 1 || result.Metrics.Exposed != 0 {
		t.Fatalf("missing lookup should fail closed without scan: %+v", result)
	}
}
