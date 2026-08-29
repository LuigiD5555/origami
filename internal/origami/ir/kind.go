package ir

// Kind identifies the semantics of an Origami representation node.
type Kind string

const (
	KindLiteral Kind = "LITERAL"
	KindRef     Kind = "REF"
	KindConcat  Kind = "CONCAT"
	KindRepeat  Kind = "REPEAT"
)

func (k Kind) Valid() bool {
	switch k {
	case KindLiteral, KindRef, KindConcat, KindRepeat:
		return true
	default:
		return false
	}
}
