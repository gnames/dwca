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

	// try to build hierarchy out of parent-child relationship
	if !a.flatHierarchy() {
		slog.Info("Building hierarchy for Core", "file", a.metaSimple.Location)
		err := a.buildHierarchy()
		if err != nil {
			return err
		}
		slog.Info("Hierarchy built", "file", a.metaSimple.Location)
	}

	// context for the whole process
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// error group and waiting group to handle concurrent processing
	g, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	// start workers
	slog.Info("Processing Core rows", "file", a.metaSimple.Location)
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

	err := a.CoreStream(ctx, chIn)
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
	p := <-a.gnpPool
	defer func() {
		a.gnpPool <- p
	}()

	for v := range chIn {
		row, err := a.processCoreRow(p, v, maxIdx)
		if err != nil {
			return err
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

func (a *arch) updateOutputCore(maxIdx int) {
	terms := []string{
		"scientificNameString",
		"canonicalFormFull",
		"canonicalForm",
		"canonicalFormStem",
		"breadcrumbNames",
		"breadcrumbRanks",
		"breadcrumbIds",
	}

	var idx int
	for i, v := range terms {
		idx = maxIdx + i + 1
		term := "https://terms.speciesfilegroup.org/" + v
		a.outputMeta.Core.Fields = append(
			a.outputMeta.Core.Fields,
			meta.Field{Term: term, Idx: idx, Index: strconv.Itoa(idx)},
		)
	}
	if _, ok := a.metaSimple.FieldsData["acceptednameusageid"]; !ok {
		idx++
		a.outputMeta.Core.Fields = append(
			a.outputMeta.Core.Fields,
			meta.Field{Term: "http://rs.tdwg.org/dwc/terms/acceptedNameUsageID",
				Idx:   idx,
				Index: strconv.Itoa(idx)},
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
	fields := meta.Headers(idx, a.outputMeta.Core.Fields)

	delim := a.outputMeta.Core.FieldsTerminatedBy
	return a.dcFile.ExportCSVStream(ctx, file, fields, delim, chOut)
}

func (a *arch) flatHierarchy() bool {
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
		nameFull := strings.TrimSpace(name + " " + author)
		nameFull, canFull, can, stem := parsedData(p, nameFull, "")
		row = append(row, nameFull, canFull, can, stem)
		res = row

	case diagn.SciNameFull, diagn.SciNameUnknown:
		name, auth := a.taxon.genNameAu(row)
		nameFull, canFull, can, stem := parsedData(p, name, auth)
		row = append(row, nameFull, canFull, can, stem)
		res = row

	case diagn.SciNameComposite:
		slog.Error("dwca.ProcessCoreRow: SciNameComposite not implemented yet")
		return nil, fmt.Errorf("not implemented yet")

	default:
		slog.Error("dwca.ProcessCoreRow: cannot process Core row")
		return nil, fmt.Errorf("cannot process Core row")
	}

	if a.flatHierarchy() {
		hNames := make([]string, len(a.taxon.hierarchy))
		ranks := make([]string, len(a.taxon.hierarchy))

		for i, v := range a.taxon.hierarchy {
			hNames[i] = row[v.index]
			ranks[i] = v.rank
		}
		breadcrumbs := strings.Join(hNames, "|")
		bcRanks := strings.Join(ranks, "|")
		res = append(res, breadcrumbs, bcRanks, "")
	} else if len(a.hierarchy) > 0 {
		taxonID := row[a.metaSimple.Index]
		ts, rs, ids := a.getBreadcrumbs(taxonID)
		res = append(res, ts, rs, ids)
	} else {
		res = append(res, "", "", "")
	}

	switch a.dgn.SynonymType {
	case diagn.SynAcceptedID:
		var taxonID, accID string
		aIdx := a.taxon.acceptedNameUsageID
		txIdx := a.taxon.taxonID
		if aIdx != -1 && txIdx != -1 {
			accID = row[aIdx]
			taxonID = row[txIdx]
			if accID == taxonID {
				res[aIdx] = ""
			}
		}
	case diagn.SynHierarchy, diagn.SynUnknown:
		if a.isSynonym(row) {
			parentID := a.parentID(row)
			res = append(res, parentID)
		} else {
			res = append(res, "")
		}
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

func parsedData(p gnparser.GNparser, name, auth string) (nameFull, canFull, can, stem string) {
	name = strings.TrimSpace(name)
	auth = strings.TrimSpace(auth)

	parsed := p.ParseName(name)
	if !parsed.Parsed {
		return name, "", "", ""
	}
	if auth != "" && parsed.Authorship == nil {
		return name + " " + auth,
			parsed.Canonical.Full,
			parsed.Canonical.Simple,
			parsed.Canonical.Stemmed
	}
	return name,
		parsed.Canonical.Full,
		parsed.Canonical.Simple,
		parsed.Canonical.Stemmed
}
