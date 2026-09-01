package temporal

import (
	"fmt"
	"sort"
	"strings"
)

const CausalTraceSchema = "origami.causal-trace.r1"

type CausalEventType string

const (
	EventSpawn    CausalEventType = "SPAWN"
	EventStart    CausalEventType = "START"
	EventComplete CausalEventType = "COMPLETE"
	EventObserve  CausalEventType = "OBSERVE"
	EventCommit   CausalEventType = "COMMIT"
	EventConflict CausalEventType = "CONFLICT"
	EventResidual CausalEventType = "RESIDUAL"
	EventCancel   CausalEventType = "CANCEL"
)

var validCausalEventTypes = map[CausalEventType]struct{}{
	EventSpawn: {}, EventStart: {}, EventComplete: {}, EventObserve: {},
	EventCommit: {}, EventConflict: {}, EventResidual: {}, EventCancel: {},
}

// CausalEvent is an immutable event-sourcing record. ParentIDs express causal
// precedence; LogicalClock is deterministic ordering metadata, not wall time.
type CausalEvent struct {
	ID              string          `json:"id"`
	Type            CausalEventType `json:"type"`
	Actor           string          `json:"actor,omitempty"`
	ParentIDs       []string        `json:"parent_ids,omitempty"`
	LogicalClock    uint64          `json:"logical_clock"`
	InputStateHash  string          `json:"input_state_hash,omitempty"`
	OutputDeltaHash string          `json:"output_delta_hash,omitempty"`
	OutputStateHash string          `json:"output_state_hash,omitempty"`
	EvidenceRefs    []string        `json:"evidence_refs,omitempty"`
}

type CausalCheckpoint struct {
	EventID   string `json:"event_id"`
	StateHash string `json:"state_hash"`
}

type CausalTrace struct {
	Schema           string             `json:"schema"`
	ID               string             `json:"id"`
	InitialStateHash string             `json:"initial_state_hash"`
	Events           []CausalEvent      `json:"events"`
	Checkpoints      []CausalCheckpoint `json:"checkpoints,omitempty"`
	FinalStateHash   string             `json:"final_state_hash,omitempty"`
}

func (t CausalTrace) Validate() error {
	if t.Schema != "" && t.Schema != CausalTraceSchema {
		return fmt.Errorf("unsupported causal trace schema %q", t.Schema)
	}
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("causal trace id is required")
	}
	if strings.TrimSpace(t.InitialStateHash) == "" {
		return fmt.Errorf("initial_state_hash is required")
	}

	byID := make(map[string]CausalEvent, len(t.Events))
	for _, event := range t.Events {
		if strings.TrimSpace(event.ID) == "" {
			return fmt.Errorf("causal event id is required")
		}
		if _, ok := validCausalEventTypes[event.Type]; !ok {
			return fmt.Errorf("unsupported causal event type %q", event.Type)
		}
		if _, exists := byID[event.ID]; exists {
			return fmt.Errorf("duplicate causal event %q", event.ID)
		}
		byID[event.ID] = event
	}

	for _, event := range t.Events {
		seenParents := map[string]struct{}{}
		for _, parentID := range event.ParentIDs {
			if parentID == event.ID {
				return fmt.Errorf("causal event %q cannot parent itself", event.ID)
			}
			if _, exists := byID[parentID]; !exists {
				return fmt.Errorf("causal event %q references unknown parent %q", event.ID, parentID)
			}
			if _, duplicate := seenParents[parentID]; duplicate {
				return fmt.Errorf("causal event %q repeats parent %q", event.ID, parentID)
			}
			seenParents[parentID] = struct{}{}
		}
	}

	if _, err := t.Frontiers(); err != nil {
		return err
	}
	for _, checkpoint := range t.Checkpoints {
		if _, exists := byID[checkpoint.EventID]; !exists {
			return fmt.Errorf("checkpoint references unknown event %q", checkpoint.EventID)
		}
		if strings.TrimSpace(checkpoint.StateHash) == "" {
			return fmt.Errorf("checkpoint %q requires state_hash", checkpoint.EventID)
		}
	}
	return nil
}

// Frontiers derives causal frames from the DAG. Events in the same frontier are
// causally independent at that depth and may overlap in wall-clock execution.
func (t CausalTrace) Frontiers() ([][]CausalEvent, error) {
	byID := make(map[string]CausalEvent, len(t.Events))
	indegree := make(map[string]int, len(t.Events))
	children := make(map[string][]string, len(t.Events))
	for _, event := range t.Events {
		byID[event.ID] = event
		indegree[event.ID] = len(event.ParentIDs)
		for _, parentID := range event.ParentIDs {
			children[parentID] = append(children[parentID], event.ID)
		}
	}

	ready := make([]string, 0)
	for id, degree := range indegree {
		if degree == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)

	frontiers := make([][]CausalEvent, 0)
	visited := 0
	for len(ready) > 0 {
		currentIDs := append([]string(nil), ready...)
		ready = nil
		frontier := make([]CausalEvent, 0, len(currentIDs))
		for _, id := range currentIDs {
			frontier = append(frontier, byID[id])
			visited++
		}
		sort.Slice(frontier, func(i, j int) bool {
			if frontier[i].LogicalClock != frontier[j].LogicalClock {
				return frontier[i].LogicalClock < frontier[j].LogicalClock
			}
			return frontier[i].ID < frontier[j].ID
		})
		frontiers = append(frontiers, frontier)

		nextSet := map[string]struct{}{}
		for _, id := range currentIDs {
			for _, childID := range children[id] {
				indegree[childID]--
				if indegree[childID] == 0 {
					nextSet[childID] = struct{}{}
				}
			}
		}
		for id := range nextSet {
			ready = append(ready, id)
		}
		sort.Strings(ready)
	}
	if visited != len(t.Events) {
		return nil, fmt.Errorf("causal trace contains a dependency cycle")
	}
	return frontiers, nil
}
