package dwca

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gnames/dwca/internal/ent/dcfile"
	"github.com/gnames/dwca/internal/ent/diagn"
	"github.com/gnames/dwca/pkg/ent/meta"
	"github.com/gnames/gnparser"
	"golang.org/x/sync/errgroup"
)

func (a *arch) processCoreOutput() error {
	chIn := make(chan []string)
	chOut := make(chan []string)

	// find the highest index in the core file
	// we can remove fields after highest index
	// and add new fields after that index
	maxIdx := slices.MaxFunc(a.meta.Core.Fields, func(a, b meta.Field) int {
		return cmp.Compare(a.Idx, b.Idx)
	}).Idx

	// taxon is a helper object to handle DarwinCore fields that are relevant
	// for name and hierarchy.
	a.taxon = a.newTaxon()

	// add new fields to Core metadata
	a.updateOutputCore(maxIdx)

	// context for the whole process
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// error group and waiting group to handle concurrent processing
	g, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	// start workers
	for i := 0; i < a.cfg.JobsNum; i++ {
		wg.Add(1)
		g.Go(func() error {
			defer wg.Done()
			return a.coreWorker(ctx, chIn, chOut, maxIdx)
		})
	}

	// close chOut when all workers are done
	go func() {
		wg.Wait()
		close(chOut)
	}()

	// save output to a file
	g.Go(func() error {
		return a.saveCoreOutput(ctx, chOut)
	})

	_, err := a.CoreStream(ctx, chIn)
	if err != nil {
		return err
	}

	if err := g.Wait(); err != nil {
		if _, ok := err.(*dcfile.ErrContext); ok {
			return err
		}
		return err
	}

	return nil
}

func (a *arch) coreWorker(
	ctx context.Context,
	chIn <-chan []string,
	chOut chan<- []string,
	maxIdx int,
) error {
	var err error
	var row []string
	p := <-a.gnpPool
	defer func() {
		a.gnpPool <- p
	}()

	for v := range chIn {
		if a.isNormalized() {
			row = v
		} else {
			row, err = a.processCoreRow(p, v, maxIdx)
			if err != nil {
				for range chIn {
				}
				return err
			}
		}

		select {
		case <-ctx.Done():
			return &dcfile.ErrContext{Err: ctx.Err()}
		default:
			chOut <- row
		}
	}
	return nil
}

func (a *arch) isNormalized() bool {
	if _, ok := a.metaSimple.FieldsData["scientificnamestring"]; ok {
		return true
	}
	return false
}

func (a *arch) updateOutputCore(maxIdx int) {
	if a.isNormalized() {
		return
	}

	terms := []string{
		"scientificNameString",
	}

	var idx int
	for i, v := range terms {
		idx = maxIdx + i + 1
		term := "https://terms.speciesfilegroup.org/" + v
		a.outputMeta.EMLFile = "eml.xml"
		a.outputMeta.Core.Fields = append(
			a.outputMeta.Core.Fields,
			meta.Field{Term: term, Idx: idx, Index: strconv.Itoa(idx)},
		)
	}
	ext := filepath.Ext(a.metaSimple.Location)
	location := a.metaSimple.Location[:len(a.metaSimple.Location)-len(ext)] + ".txt"

	delim := ","
	if a.cfg.OutputCSVType == "tsv" {
		delim = `\t`
	}

	a.outputMeta.Core.Files.Location = location
	a.outputMeta.Core.FieldsEnclosedBy = `"`
	a.outputMeta.Core.FieldsTerminatedBy = delim
	a.outputMeta.Core.IgnoreHeaderLines = "1"
	a.outputMeta.Core.LinesTerminatedBy = `\n`
}

func (a *arch) saveCoreOutput(ctx context.Context, chOut <-chan []string) error {
	file := a.outputMeta.Core.Files.Location

	idx := a.outputMeta.Core.ID.Idx
	headers := meta.Headers(idx, a.outputMeta.Core.Fields)

	delim := a.outputMeta.Core.FieldsTerminatedBy
	return a.dcFile.ExportCSVStream(ctx, file, headers, delim, chOut)
}

func (a *arch) flatHierarchy() bool {
	// if tree hierarchy exists, favor it over flat one.
	if a.taxon.higherTaxonID != -1 ||
		a.taxon.parentNameUsageID != -1 {
		return false
	}
	return len(a.taxon.hierarchy) > 1
}

func (a *arch) processCoreRow(
	p gnparser.GNparser,
	row []string,
	maxIdx int,
) ([]string, error) {
	var res []string
	// add empty fields if row is short, cut fields larger than maxIdx.
	row = a.normalizeRow(row, maxIdx)
	switch a.dgn.SciNameType {

	case diagn.SciNameCanonical:
		name, author := a.taxon.genNameAu(row)
		nameStr := strings.TrimSpace(name + " " + author)
		nameStr = getFullName(p, nameStr, "")
		row = append(row, nameStr)
		res = row

	case diagn.SciNameFull, diagn.SciNameUnknown:
		name, auth := a.taxon.genNameAu(row)
		nameStr := getFullName(p, name, auth)
		row = append(row, nameStr)
		res = row

	case diagn.SciNameComposite:
		slog.Error("dwca.ProcessCoreRow: SciNameComposite not implemented yet")
		return nil, fmt.Errorf("not implemented yet")

	default:
		slog.Error("dwca.ProcessCoreRow: cannot process Core row")
		return nil, fmt.Errorf("cannot process Core row")
	}

	return res, nil
}

func (a *arch) parentID(row []string) string {
	pIdx := a.taxon.parentNameUsageID
	if pIdx != -1 {
		return row[pIdx]
	}
	pIdx = a.taxon.higherTaxonID
	if pIdx != -1 {
		return row[pIdx]
	}
	return ""
}

func (a *arch) isSynonym(row []string) bool {
	syn := []string{"synonym", "homonym", "misapplied", "ambiguous"}
	synPart := []string{"synonym", "miss", "un"}
	tsIdx := a.taxon.taxonomicStatus

	if tsIdx == -1 {
		return false
	}

	st := strings.TrimSpace(row[tsIdx])
	for i := range syn {
		if st == syn[i] {
			return true
		}
	}
	for i := range synPart {
		if strings.Contains(st, synPart[i]) {
			return true
		}
	}
	return false
}

func (a *arch) normalizeRow(row []string, maxIdx int) []string {
	l := len(row)
	for _, v := range a.meta.Core.Fields {
		if v.Idx == -1 {
			continue
		}
		if v.Idx >= l {
			row = append(row, "")
		}
	}
	return row[0 : maxIdx+1]
}

// getFullName checks if a scientificName field contains name with authorship.
// if yes, it does not try to append it with authorship from
// scientificNameAuthorship field.
func getFullName(p gnparser.GNparser, name, auth string) string {
	name = strings.TrimSpace(name)
	auth = strings.TrimSpace(auth)

	parsed := p.ParseName(name)
	if !parsed.Parsed {
		return name
	}
	if auth != "" && parsed.Authorship == nil {
		return name + " " + auth
	}
	return name
}
