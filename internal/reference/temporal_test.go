package reference

import "testing"

func temporalFixture(t *testing.T) CoreRun{
	t.Helper()
	initial:=CoreState{"signal":{Status:StatusPresent,Value:"low"}}
	rules:=[]Rule{
		{ID:"R-rise",When:[]Predicate{{Source:"context",Key:"phase",Equals:"rise"}},Then:[]Effect{{Entity:"signal",Status:StatusPresent,Value:"high"}}},
		{ID:"R-fall",When:[]Predicate{{Source:"context",Key:"phase",Equals:"fall"}},Then:[]Effect{{Entity:"signal",Status:StatusPresent,Value:"low"}}},
	}
	run,err:=ExecuteCore(initial,[]Context{{ID:"C1",Values:map[string]string{"phase":"rise"}},{ID:"C2",Values:map[string]string{"phase":"fall"}}},rules);if err!=nil{t.Fatal(err)};return run
}

func TestADV008PropertyExistsOnlyAcrossTrajectory(t *testing.T){
	run:=temporalFixture(t)
	contract:=TemporalObservationContract{ID:"OBS-pulse",Authority:"temporal-test",Predicate:TemporalPredicate{Entity:"signal",Relation:TemporalSequence,Values:[]string{"low","high","low"}}}
	before,err:=TraceHash(run.Trace);if err!=nil{t.Fatal(err)}
	obs,err:=ObserveTemporal(run,contract);if err!=nil{t.Fatal(err)}
	if !obs.Matched{t.Fatalf("ADV-008: trajectory property not detected: %v",obs.ObservedValues)}
	if obs.TraceHash!=before{t.Fatal("ADV-008: temporal observation detached from immutable trace")}
	for i,state:=range run.Trace{if v,ok:=state["pulse"];ok&&v.Status==StatusPresent{t.Fatalf("ADV-008: final property was smuggled into frame %d: %+v",i,v)}}
}

func TestADV008NoSingleFrameCanSatisfyTemporalSequence(t *testing.T){
	run:=temporalFixture(t);contract:=TemporalObservationContract{ID:"OBS-pulse",Authority:"temporal-test",Predicate:TemporalPredicate{Entity:"signal",Relation:TemporalSequence,Values:[]string{"low","high","low"}}}
	for i,state:=range run.Trace{single:=CoreRun{Trace:[]CoreState{state}};obs,err:=ObserveTemporal(single,contract);if err!=nil{t.Fatal(err)};if obs.Matched{t.Fatalf("ADV-008: single frame %d incorrectly contained trajectory-only property",i)}}
}

func TestADV008ReorderingTrajectoryDestroysProperty(t *testing.T){
	run:=temporalFixture(t);contract:=TemporalObservationContract{ID:"OBS-pulse",Authority:"temporal-test",Predicate:TemporalPredicate{Entity:"signal",Relation:TemporalSequence,Values:[]string{"low","high","low"}}}
	reordered:=CoreRun{Trace:[]CoreState{run.Trace[1],run.Trace[0],run.Trace[2]}}
	obs,err:=ObserveTemporal(reordered,contract);if err!=nil{t.Fatal(err)};if obs.Matched{t.Fatalf("ADV-008: reordered trajectory preserved order-dependent property: %v",obs.ObservedValues)}
}

func TestADV008TemporalObserverDoesNotWriteDerivedAnswer(t *testing.T){
	run:=temporalFixture(t);before,err:=TraceHash(run.Trace);if err!=nil{t.Fatal(err)}
	_,err=ObserveTemporal(run,TemporalObservationContract{ID:"OBS-change",Authority:"temporal-test",Predicate:TemporalPredicate{Entity:"signal",Relation:TemporalChange}});if err!=nil{t.Fatal(err)}
	after,err:=TraceHash(run.Trace);if err!=nil{t.Fatal(err)};if before!=after{t.Fatal("ADV-008: observer mutated trajectory to encode derived answer")}
}
