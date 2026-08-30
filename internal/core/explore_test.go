package core

import "testing"

func TestCombinatorialOracleEnumeratesFullDeclaredProduct(t *testing.T){
	domain:=[]Value{{Status:StatusPresent,Data:"0"},{Status:StatusPresent,Data:"1"},{Status:StatusUnknown}}
	states,err:=EnumerateStates([]string{"C","A","B"},domain);if err!=nil{t.Fatal(err)}
	if len(states)!=27{t.Fatalf("expected 3^3=27 initial states, got %d",len(states))}
	seen:=map[string]struct{}{};for _,state:=range states{seen[StateHash(state)]=struct{}{}};if len(seen)!=27{t.Fatalf("enumeration collapsed distinct states: %d",len(seen))}
}

func TestDeclaredContradictionTerminatesBeforeMutation(t *testing.T){
	machine:=NewMachine([]Rule{
		{ID:"contradiction",When:[]Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"},{Scope:"state",Key:"B",Op:"data_eq",Data:"inactive"}},Contradiction:&Contradiction{Reason:"A requires B but B is inactive"}},
		{ID:"would-mutate",When:[]Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"}},Then:[]Effect{{Key:"B",Mode:"set",Status:StatusPresent,Data:"active"}}},
	},4)
	initial:=State{Values:map[string]Value{"A":{Status:StatusPresent,Data:"active"},"B":{Status:StatusPresent,Data:"inactive"}}}
	trajectory,err:=Execute(initial,machine,nil);if err!=nil{t.Fatal(err)}
	if trajectory.Terminated!="CONTRADICTION"{t.Fatalf("expected contradiction, got %s",trajectory.Terminated)}
	if trajectory.States[len(trajectory.States)-1].Values["B"].Data!="inactive"{t.Fatal("contradiction rule must terminate before mutation")}
	if len(trajectory.Evidence[0].Contradictions)!=1||trajectory.Evidence[0].Contradictions[0].RuleID!="contradiction"{t.Fatalf("contradiction provenance missing: %+v",trajectory.Evidence[0])}
}

func TestExploreReportsEveryInitialStateWithoutSilentPruning(t *testing.T){
	machine:=NewMachine([]Rule{{ID:"copy",When:[]Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"1"}},Then:[]Effect{{Key:"B",Mode:"set",Status:StatusPresent,Data:"1"}}}},3)
	report,err:=Explore([]string{"A","B"},[]Value{{Status:StatusPresent,Data:"0"},{Status:StatusPresent,Data:"1"}},machine,nil);if err!=nil{t.Fatal(err)}
	if report.Initials!=4||len(report.Trajectories)!=4{t.Fatalf("exploration silently pruned state space: %+v",report)}
}
