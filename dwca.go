package dwca

type factory struct {
	url string
}

func New() Factory {
	return &factory{}
}

func (f *factory) Fetch(url string) error {
	return nil
}
