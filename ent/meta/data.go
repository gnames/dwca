package meta

type MetaSimple struct {
	CoreData
	ExtensionsData map[string]ExtensionData
}

type CoreData struct {
	Index      int
	Term       string
	TermFull   string
	Location   string
	FieldsData map[string]FieldData
	FieldsIdx  map[int]FieldData
}

type ExtensionData struct {
	CoreIndex  int
	Location   string
	FieldsData map[string]FieldData
	FieldsIdx  map[int]FieldData
}

type FieldData struct {
	Index    int
	Term     string
	TermFull string
}
