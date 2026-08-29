package reference

import "testing"

func TestDependencyFoldPreservesGuardsAndProof(t *testing.T) {
	steps := []FoldStep{
		{From:"A", To:"B", Rule:"r1", Guards:[]FoldGuard{{Kind:"condition",Key:"mode",Value:"x"}}},
		{From:"B", To:"C", Rule:"r2", Guards:[]FoldGuard{{Kind:"exception",Key:"blocked",Value:"false"}}},
		{From:"C", To:"D", Rule:"r3"},
	}
	edge, err := FoldDependency(steps); if err != nil { t.Fatal(err) }
	if edge.From!="A" || edge.To!="D" || edge.Depth!=3 { t.Fatalf("bad edge: %+v", edge) }
	if len(edge.Guards)!=2 { t.Fatalf("guards lost: %+v", edge.Guards) }
	if !VerifyFold(edge) { t.Fatal("valid fold did not verify") }
}

func TestDependencyFoldRejectsBrokenPath(t *testing.T) {
	_, err := FoldDependency([]FoldStep{{From:"A",To:"B"},{From:"C",To:"D"}})
	if err == nil { t.Fatal("expected non-contiguous path rejection") }
}

func TestDependencyFoldDetectsCanonicalMutation(t *testing.T) {
	edge, err := FoldDependency([]FoldStep{{From:"A",To:"B"},{From:"B",To:"C"}}); if err != nil { t.Fatal(err) }
	edge.Canonical[1].To = "X"
	if VerifyFold(edge) { t.Fatal("mutated canonical path verified") }
}

func TestFoldLevelBoundsExposure(t *testing.T) {
	steps:=make([]FoldStep, 1000)
	for i:=range steps { steps[i]=FoldStep{From:node(i),To:node(i+1),Rule:"next"} }
	level,err:=FoldLevel(steps,64); if err!=nil { t.Fatal(err) }
	if len(level)!=16 { t.Fatalf("want 16 blocks, got %d",len(level)) }
	for _,edge:=range level { if !VerifyFold(edge) { t.Fatal("block proof failed") } }
}

func node(i int) string {
	const digits="0123456789"
	if i==0{return "N0"}
	b:=make([]byte,0,8); for i>0 { b=append(b,digits[i%10]); i/=10 }
	for l,r:=0,len(b)-1;l<r;l,r=l+1,r-1 { b[l],b[r]=b[r],b[l] }
	return "N"+string(b)
}
