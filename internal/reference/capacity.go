package reference

import "fmt"

type CapacityFailure string

const (
	CapacityOK CapacityFailure = "OK"
	FailureWindowExceeded CapacityFailure = "WINDOW_EXCEEDED"
	FailureClosureExplosion CapacityFailure = "CLOSURE_EXPLOSION"
	FailureRouteIntegrity CapacityFailure = "ROUTE_INTEGRITY_FAILURE"
)

type WorkingBudget struct {
	TokenEquivalent int `json:"token_equivalent"`
	MaxActiveBytes int `json:"max_active_bytes"`
	MaxSemanticUnits int `json:"max_semantic_units"`
}

type AccessMeasurement struct {
	TotalSemanticUnits int `json:"total_semantic_units"`
	RequestedAddresses int `json:"requested_addresses"`
	VisitedAddresses int `json:"visited_addresses"`
	UnfoldedSemanticUnits int `json:"unfolded_semantic_units"`
	UnrelatedUnfolded int `json:"unrelated_unfolded"`
	ActiveBytes int `json:"active_bytes"`
	TokenEquivalentEstimate int `json:"token_equivalent_estimate"`
	AccessFraction float64 `json:"access_fraction"`
	WithinBudget bool `json:"within_budget"`
	Failure CapacityFailure `json:"failure"`
}

// MeasureAccess evaluates the active semantic working set. TokenEquivalent is a
// reporting/reference budget only; Origami's internal state is not defined in tokens.
func MeasureAccess(totalUnits,requested,visited,unfolded,unrelated,activeBytes int,budget WorkingBudget)(AccessMeasurement,error){
	if totalUnits<=0{return AccessMeasurement{},fmt.Errorf("total semantic units must be positive")};if requested<0||visited<0||unfolded<0||unrelated<0||activeBytes<0{return AccessMeasurement{},fmt.Errorf("capacity measurements cannot be negative")};if unfolded>totalUnits{return AccessMeasurement{},fmt.Errorf("unfolded units exceed represented units")}
	// Four bytes/token is deliberately an estimate for comparability with model
	// windows, never an Origami semantic unit definition.
	tokenEstimate:=(activeBytes+3)/4
	within:=true;failure:=CapacityOK
	if budget.MaxSemanticUnits>0&&unfolded>budget.MaxSemanticUnits{within=false;failure=FailureClosureExplosion}
	if budget.MaxActiveBytes>0&&activeBytes>budget.MaxActiveBytes{within=false;failure=FailureWindowExceeded}
	if budget.TokenEquivalent>0&&tokenEstimate>budget.TokenEquivalent{within=false;failure=FailureWindowExceeded}
	return AccessMeasurement{TotalSemanticUnits:totalUnits,RequestedAddresses:requested,VisitedAddresses:visited,UnfoldedSemanticUnits:unfolded,UnrelatedUnfolded:unrelated,ActiveBytes:activeBytes,TokenEquivalentEstimate:tokenEstimate,AccessFraction:float64(unfolded)/float64(totalUnits),WithinBudget:within,Failure:failure},nil
}
