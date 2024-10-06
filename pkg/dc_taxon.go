package dwca

import (
	"cmp"
	"slices"
	"strings"
)

type taxon struct {
	scientificName,
	scientificNameAuthorship,
	domain,
	kingdom,
	phylum,
	class,
	order,
	superfamily,
	family,
	subfamily,
	tribe,
	subtribe,
	genericName,
	genus,
	subgenus,
	infragenericEpithet,
	specificEpithet,
	infraspecificEpithet,
	scientificNameRank,
	taxonRank,
	taxonomicStatus,
	acceptedNameUsageID,
	higherTaxonID,
	parentNameUsageID,
	taxonID int
	hierarchy []taxonHierarchy
}

type taxonHierarchy struct {
	sortBy int
	rank   string
	index  int
}

func (a *arch) newTaxon() *taxon {
	res := taxon{
		scientificName:           -1,
		scientificNameAuthorship: -1,
		domain:                   -1,
		kingdom:                  -1,
		phylum:                   -1,
		class:                    -1,
		order:                    -1,
		superfamily:              -1,
		family:                   -1,
		subfamily:                -1,
		tribe:                    -1,
		subtribe:                 -1,
		genericName:              -1,
		genus:                    -1,
		subgenus:                 -1,
		infragenericEpithet:      -1,
		specificEpithet:          -1,
		infraspecificEpithet:     -1,
		scientificNameRank:       -1,
		taxonRank:                -1,
		taxonomicStatus:          -1,
		acceptedNameUsageID:      -1,
		higherTaxonID:            -1,
		parentNameUsageID:        -1,
		taxonID:                  -1,
	}
	cr := a.metaSimple.CoreData
	for k, v := range cr.FieldsData {
		switch k {
		case "scientificname":
			res.scientificName = v.Index
		case "scientificnameauthorship":
			res.scientificNameAuthorship = v.Index
		case "domain":
			res.domain = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 5, rank: "domain", index: v.Index})
		case "kingdom":
			res.kingdom = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 10, rank: "kingdom", index: v.Index})
		case "phylum":
			res.phylum = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 20, rank: "phylum", index: v.Index})
		case "class":
			res.class = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 30, rank: "class", index: v.Index})
		case "order":
			res.order = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 40, rank: "order", index: v.Index})
		case "superfamily":
			res.superfamily = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 50, rank: "superfamily", index: v.Index})
		case "family":
			res.family = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 60, rank: "family", index: v.Index})
		case "subfamily":
			res.subfamily = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 70, rank: "subfamily", index: v.Index})
		case "tribe":
			res.tribe = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 80, rank: "tribe", index: v.Index})
		case "subtribe":
			res.subtribe = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 90, rank: "subtribe", index: v.Index})
		case "genericname":
			res.genericName = v.Index
		case "genus":
			res.genus = v.Index
			res.hierarchy = append(res.hierarchy,
				taxonHierarchy{sortBy: 100, rank: "genus", index: v.Index})
		case "subgenus":
			res.subgenus = v.Index
		case "infragenericepithet":
			res.infragenericEpithet = v.Index
		case "specificepithet":
			res.specificEpithet = v.Index
		case "infraspecificepithet":
			res.infraspecificEpithet = v.Index
		case "scientificnamerank":
			res.scientificNameRank = v.Index
		case "taxonrank":
			res.taxonRank = v.Index
		case "taxonomicstatus":
			res.taxonomicStatus = v.Index
		case "acceptednameusageid":
			res.acceptedNameUsageID = v.Index
		case "highertaxonid":
			res.higherTaxonID = v.Index
		case "parentnameusageid":
			res.parentNameUsageID = v.Index
		case "taxonid":
			res.taxonID = v.Index
		}
	}
	slices.SortFunc(res.hierarchy, func(a, b taxonHierarchy) int {
		return cmp.Compare(a.sortBy, b.sortBy)
	})

	return &res
}

func (n *taxon) genNameAu(row []string) (string, string) {
	if n.scientificName == -1 {
		return "", ""
	}
	sn := row[n.scientificName]
	var aus string
	if n.scientificNameAuthorship != -1 {
		aus = row[n.scientificNameAuthorship]
	}
	sn = strings.TrimSpace(sn)
	aus = strings.TrimSpace(aus)
	return sn, aus
}
