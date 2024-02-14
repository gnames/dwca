package diagn

type SynonymType int

const (
	SynUnknown SynonymType = iota
	SynAcceptedID
	SynHierarchy
	SynExtension
)

func (st SynonymType) String() string {
	switch st {
	case SynAcceptedID:
		return "accepted ID"
	case SynHierarchy:
		return "hierarchy"
	case SynExtension:
		return "extension"
	default:
		return "unknown"
	}
}
