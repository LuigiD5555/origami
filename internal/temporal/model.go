package temporal

import "github.com/LuigiD5555/origami/internal/automaton"

const Schema = "origami.temporal-program.r0"

type Program struct {
	Schema          string              `json:"schema"`
	ID              string              `json:"id"`
	Automaton       automaton.Automaton `json:"automaton"`
	InitialState    automaton.State     `json:"initial_state,omitempty"`
	MaxSteps        int                 `json:"max_steps"`
	CheckpointEvery int                 `json:"checkpoint_every,omitempty"`
}

type Frame struct {
	Step    int                `json:"step"`
	Changes []automaton.Change `json:"changes"`
	Stable  bool               `json:"stable,omitempty"`
}

type Checkpoint struct {
	Step  int             `json:"step"`
	State automaton.State `json:"state"`
}

type Trace struct {
	Schema      string          `json:"schema"`
	ProgramID   string          `json:"program_id"`
	Initial     automaton.State `json:"initial"`
	Frames      []Frame         `json:"frames"`
	Checkpoints []Checkpoint    `json:"checkpoints,omitempty"`
	Final       automaton.State `json:"final"`
	StableAt    int             `json:"stable_at,omitempty"`
}

type Metrics struct {
	Steps                int     `json:"steps"`
	Changes              int     `json:"changes"`
	Checkpoints          int     `json:"checkpoints"`
	ExplicitSnapshotByte int     `json:"explicit_snapshot_bytes"`
	FoldedTraceByte      int     `json:"folded_trace_bytes"`
	TemporalFoldRatio    float64 `json:"temporal_fold_ratio"`
}
