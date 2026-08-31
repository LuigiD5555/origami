package temporal

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/LuigiD5555/origami/internal/automaton"
)

func (p Program) Validate() error {
	if p.Schema != "" && p.Schema != Schema {
		return fmt.Errorf("unsupported temporal schema %q", p.Schema)
	}
	if p.ID == "" {
		return fmt.Errorf("program id is required")
	}
	if p.MaxSteps <= 0 {
		return fmt.Errorf("max_steps must be > 0")
	}
	if p.CheckpointEvery < 0 {
		return fmt.Errorf("checkpoint_every must be >= 0")
	}
	return p.Automaton.Validate()
}

func initialState(p Program) automaton.State {
	if len(p.InitialState) > 0 {
		return automaton.CloneState(p.InitialState)
	}
	return p.Automaton.Initial()
}

func Run(p Program) (Trace, error) {
	if err := p.Validate(); err != nil {
		return Trace{}, err
	}
	state := initialState(p)
	trace := Trace{Schema: Schema, ProgramID: p.ID, Initial: automaton.CloneState(state)}
	if p.CheckpointEvery > 0 {
		trace.Checkpoints = append(trace.Checkpoints, Checkpoint{Step: 0, State: automaton.CloneState(state)})
	}
	for step := 1; step <= p.MaxSteps; step++ {
		next, result, err := automaton.Step(p.Automaton, state, step)
		if err != nil {
			return Trace{}, err
		}
		trace.Frames = append(trace.Frames, Frame{Step: step, Changes: result.Changes, Stable: result.Stable})
		state = next
		if p.CheckpointEvery > 0 && step%p.CheckpointEvery == 0 {
			trace.Checkpoints = append(trace.Checkpoints, Checkpoint{Step: step, State: automaton.CloneState(state)})
		}
		if result.Stable {
			trace.StableAt = step
			break
		}
	}
	trace.Final = automaton.CloneState(state)
	return trace, nil
}

func applyChanges(state automaton.State, changes []automaton.Change) (automaton.State, error) {
	out := automaton.CloneState(state)
	for _, c := range changes {
		cur, ok := out[c.CellID]
		if !ok {
			return nil, fmt.Errorf("delta references unknown cell %q", c.CellID)
		}
		if cur != c.From {
			return nil, fmt.Errorf("delta mismatch for %q: have %q want %q", c.CellID, cur, c.From)
		}
		out[c.CellID] = c.To
	}
	return out, nil
}

// Seek reconstructs one state from the closest checkpoint not after target.
// It never needs to replay frames before that checkpoint.
func Seek(trace Trace, target int) (automaton.State, error) {
	if target < 0 {
		return nil, fmt.Errorf("target step must be >= 0")
	}
	state := automaton.CloneState(trace.Initial)
	start := 0
	for _, cp := range trace.Checkpoints {
		if cp.Step <= target && cp.Step >= start {
			state = automaton.CloneState(cp.State)
			start = cp.Step
		}
	}
	frames := append([]Frame(nil), trace.Frames...)
	sort.Slice(frames, func(i, j int) bool { return frames[i].Step < frames[j].Step })
	for _, f := range frames {
		if f.Step <= start {
			continue
		}
		if f.Step > target {
			break
		}
		var err error
		state, err = applyChanges(state, f.Changes)
		if err != nil {
			return nil, err
		}
	}
	return state, nil
}

func VerifyReplay(p Program, trace Trace) error {
	actual, err := Run(p)
	if err != nil {
		return err
	}
	if !automaton.EqualState(actual.Final, trace.Final) {
		return fmt.Errorf("final state mismatch")
	}
	if len(actual.Frames) != len(trace.Frames) {
		return fmt.Errorf("frame count mismatch: %d != %d", len(actual.Frames), len(trace.Frames))
	}
	for i := range actual.Frames {
		a, _ := json.Marshal(actual.Frames[i])
		b, _ := json.Marshal(trace.Frames[i])
		if string(a) != string(b) {
			return fmt.Errorf("frame %d mismatch", i)
		}
	}
	return nil
}

func Measure(trace Trace) Metrics {
	changes := 0
	states := []automaton.State{automaton.CloneState(trace.Initial)}
	state := automaton.CloneState(trace.Initial)
	for _, f := range trace.Frames {
		changes += len(f.Changes)
		state, _ = applyChanges(state, f.Changes)
		states = append(states, automaton.CloneState(state))
	}
	explicit, _ := json.Marshal(states)
	folded, _ := json.Marshal(trace)
	ratio := 0.0
	if len(folded) > 0 {
		ratio = float64(len(explicit)) / float64(len(folded))
	}
	return Metrics{
		Steps: len(trace.Frames), Changes: changes, Checkpoints: len(trace.Checkpoints),
		ExplicitSnapshotByte: len(explicit), FoldedTraceByte: len(folded), TemporalFoldRatio: ratio,
	}
}
