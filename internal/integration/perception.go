package integration

import (
	"encoding/json"
	"fmt"

	"github.com/LuigiD5555/origami/internal/lab/observation"
	"github.com/LuigiD5555/origami/internal/memory"
)

// PerceptionNode projects an OHF/native observation into the shared addressable
// memory plane. Validation of an observation proves only that the perception
// record is well-formed and bound to a model trial; it does not promote its
// interpretation to verified semantic evidence.
func PerceptionNode(carrierID string, envelope observation.Envelope) (memory.Node, error) {
	if err := envelope.Validate(); err != nil { return memory.Node{}, err }
	body,err:=json.Marshal(envelope);if err!=nil{return memory.Node{},err}
	address:=memory.CanonicalAddress(carrierID,"perception",fmt.Sprintf("%s-%d",cleanID(envelope.Model),envelope.Trial))
	return memory.Node{
		Address:address,CarrierID:carrierID,ClusterID:"perception",Kind:"perception-observation",SemanticClass:string(envelope.Status),
		Label:fmt.Sprintf("%s trial %d perception %s",envelope.Model,envelope.Trial,envelope.Status),
		Levels:map[memory.Fidelity]memory.ContentLevel{
			memory.FidelityAbstract:{Text:fmt.Sprintf("perception %s from model trial",envelope.Status)},
			memory.FidelityEvidence:{Text:string(body),SourceRef:envelope.BundleID,Verified:false},
		},
		Provenance:memory.Provenance{SourceID:"ohf:native-observation",SourceRef:envelope.BundleID,SourceHash:memory.ContentID(string(body))},
		VerificationRef:"PERCEPTION_ONLY_NOT_SEMANTIC_AUTHORITY",
	},nil
}
