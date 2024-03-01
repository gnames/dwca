package meta

import (
	"cmp"
	"path/filepath"
	"slices"
	"strconv"
)

// Headers return headers for the output file. It takes idx parameter
// which corresponds to designated Index field of DwCA star schema.
// It also takes fields, which is a slice of fields that corresponds
// to the DwCA star schema.
func Headers(idx int, fields []Field) []string {
	lastField := slices.MaxFunc(fields, func(a, b Field) int {
		return cmp.Compare(a.Idx, b.Idx)
	})

	fieldMap := make(map[int]Field)
	for _, f := range fields {
		fieldMap[f.Idx] = f
	}
	var unknownCount int
	res := make([]string, lastField.Idx+1)
	for i := range lastField.Idx + 1 {
		if f, ok := fieldMap[i]; ok {
			term := filepath.Base(f.Term)
			res[i] = term
			continue
		}

		// Index does not have a field description, we try our best guess here.
		if i == idx {
			// it might bite later.
			res[i] = "taxonID"
			continue
		}
		unknownCount++

		term := "unknown" + strconv.Itoa(unknownCount)
		res[i] = term
	}

	return res
}
