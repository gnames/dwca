package dwca

type Factory interface {
	Fetch(url string) error
}
