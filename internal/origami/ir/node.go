package ir

// Node is the minimal R0 generative representation unit. Fields not used by
// a node kind remain zero-valued. IDs are stable within a Representation.
type Node struct {
	ID       string
	Kind     Kind
	Literal  []byte
	Target   string
	Children []string
	Count    uint64
}

// Representation is an addressable collection of nodes with one root.
type Representation struct {
	Root  string
	Nodes map[string]Node
}
