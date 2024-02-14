package diagn

type SciNameType int

const (
	SciNameUnknown SciNameType = iota
	SciNameFull
	SciNameCanonical
	SciNameComposite
)

func (s SciNameType) String() string {
	switch s {
	case SciNameFull:
		return "full"
	case SciNameCanonical:
		return "canonical"
	case SciNameComposite:
		return "composite"
	default:
		return "unknown"
	}
}
