package meta

type Data struct {
	CoreData
	ExtensionsData map[string]ExtensionData
}

type CoreData struct {
	Index      int
	Term       string
	TermFull   string
	FieldsData map[string]FieldData
}

type ExtensionData struct {
	CoreIndex  int
	FieldsData map[string]FieldData
}

type FieldData struct {
	Index    int
	Term     string
	TermFull string
}
