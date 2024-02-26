package dwca

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/gnames/gnparser"
	"golang.org/x/sync/errgroup"
)

type hNode struct {
	id       string
	parentID string
	name     string
	rank     string
}

func (a *arch) buildHierarchy() error {
	// if scientific name is not given, we cannot build hierarchy
	if _, ok := a.metaSimple.FieldsData["scientificname"]; !ok {
		return nil
	}

	// if parent is not given, we cannot build hierarchy
	if _, ok := a.metaSimple.FieldsData["parentnameusageid"]; !ok {
		if _, ok := a.metaSimple.FieldsData["highertaxonid"]; !ok {
			return nil
		}
	}

	chIn := make(chan []string)
	chOut := make(chan *hNode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	for i := 0; i < a.cfg.JobsNum; i++ {
		wg.Add(1)
		g.Go(func() error {
			defer wg.Done()
			return a.hierarchyWorker(ctx, chIn, chOut)
		})
	}

	g.Go(func() error {
		return a.createHierarchy(ctx, chOut)
	})

	// close chOut when all workers are done
	go func() {
		wg.Wait()
		close(chOut)
	}()

	err := a.CoreStream(ctx, chIn)
	if err != nil {
		return err
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (a *arch) hierarchyWorker(
	ctx context.Context,
	chIn <-chan []string,
	chOut chan<- *hNode,
) error {
	p := <-a.gnpPool
	defer func() {
		a.gnpPool <- p
	}()

	for v := range chIn {
		row, err := a.processHierarchyRow(p, v)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			chOut <- row
		}
	}

	return nil
}

func (a *arch) processHierarchyRow(p gnparser.GNparser, v []string) (*hNode, error) {
	var canonical string
	if field, ok := a.metaSimple.FieldsData["scientificname"]; ok {
		name := v[field.Index]
		canonical = name
		_, _, can, _ := parsedData(p, name, "")
		if can != "" {
			canonical = can
		}
	}

	parentID := v[a.metaSimple.FieldsData["parentnameusageid"].Index]
	if parentID == "" {
		parentID = v[a.metaSimple.FieldsData["highertaxonid"].Index]
	}

	var rank string
	if field, ok := a.metaSimple.FieldsData["taxonrank"]; ok {
		rank = v[field.Index]
	}

	res := &hNode{
		id:       v[a.metaSimple.Index],
		rank:     rank,
		name:     canonical,
		parentID: parentID,
	}

	return res, nil
}

func (a *arch) createHierarchy(ctx context.Context, chOut <-chan *hNode) error {
	for v := range chOut {
		if v.id == "" {
			continue
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			a.hierarchy[v.id] = v
		}
	}
	return nil
}

func (a *arch) getBreadcrumbs(id string) (bcTx, bcRnk, bcIdx string) {
	nodes := a.breadcrumbsNodes(id)

	ts := make([]string, len(nodes))
	rs := make([]string, len(nodes))
	is := make([]string, len(nodes))

	for i := range nodes {
		ts[i] = nodes[i].name
		rs[i] = nodes[i].rank
		is[i] = nodes[i].id
	}

	return strings.Join(ts, "|"), strings.Join(rs, "|"), strings.Join(is, "|")
}

func (a *arch) breadcrumbsNodes(id string) []*hNode {
	var res []*hNode
	var node *hNode
	var ok bool
	var currID, prevID string

	currID = id
	for currID != "" {
		if node, ok = a.hierarchy[currID]; !ok {
			slog.Warn("Hierarchy node not found", "id", currID)
			return res
		}
		res = append([]*hNode{node}, res...)
		prevID = currID
		currID = strings.TrimSpace(node.parentID)
		if currID == prevID {
			return res
		}
	}
	return res
}
