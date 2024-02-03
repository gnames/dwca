package dwca

// Archive is an interface for Darwin Core Archive objects.
type Archive interface {
	// Extract unpacks the DwCA files into a working temporary directory.
	Extract() error
}
