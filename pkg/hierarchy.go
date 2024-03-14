package dwca

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"

	"github.com/gnames/gnparser"
	"golang.org/x/sync/errgroup"
)

type hNode struct {
	id         string
	parentID   string
	acceptedID string
	name       string
	rank       string
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

	_, err := a.CoreStream(ctx, chIn)
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
			for _ = range chIn {
			}
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

	id := v[a.metaSimple.Index]
	if id == "" {
		return nil, errors.New("ID is empty")
	}

	var parentID string
	if field, ok := a.metaSimple.FieldsData["parentnameusageid"]; ok {
		parentID = v[field.Index]
	}
	if parentID == "" {
		if field, ok := a.metaSimple.FieldsData["highertaxonid"]; ok {
			parentID = v[field.Index]
		}
	}

	if parentID == id {
		parentID = ""
	}

	var acceptedID string
	if field, ok := a.metaSimple.FieldsData["acceptednameusageid"]; ok {
		acceptedID = v[field.Index]
	}

	if acceptedID == id {
		acceptedID = ""
	}

	var rank string
	if field, ok := a.metaSimple.FieldsData["taxonrank"]; ok {
		rank = v[field.Index]
	}

	res := &hNode{
		id:         strings.TrimSpace(id),
		rank:       strings.TrimSpace(rank),
		name:       strings.TrimSpace(canonical),
		parentID:   strings.TrimSpace(parentID),
		acceptedID: strings.TrimSpace(acceptedID),
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
	id = strings.TrimSpace(id)
	var res []*hNode
	var node *hNode
	var ok bool
	var currID, prevID string

	currID = id
	for currID != "" {
		if node, ok = a.hierarchy[currID]; !ok {
			slog.Warn("Hierarchy node not found, making short breadcumbs", "id", currID)
			return res
		}

		if node.parentID == "" && node.acceptedID != "" {
			currID = node.acceptedID
			continue
		}

		res = append([]*hNode{node}, res...)
		prevID = currID
		currID = node.parentID
		if currID == prevID {
			return res
		}
	}
	return res
}
