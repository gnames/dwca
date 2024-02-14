package diagn

type HierType int

const (
	HierUnknown HierType = iota
	HierTree
	HierFlat
)

func (h HierType) String() string {
	switch h {
	case HierTree:
		return "tree"
	case HierFlat:
		return "flat"
	default:
		return "unknown"
	}
}
