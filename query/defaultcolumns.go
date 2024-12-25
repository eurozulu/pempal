package query

import (
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"strings"
)

var defaultFields = map[model.ResourceType][]string{
	model.UnknownResourceType: {"type", "id", "filename"},
	model.Certificate:         {"type", "id", "subject", "issuer", "not-after", "is-ca", "filename"},
	model.PrivateKey:          {"type", "id", "public-key-algorithm", "filename"},
	model.PublicKey:           {"type", "id", "public-key-algorithm", "filename"},
	model.CertificateRequest:  {"type", "id", "subject", "version", "filename"},
	model.RevokationList:      {"type", "issuer", "number", "this-update", "filename"},
}

func DefaultColumnNames(types ...model.ResourceType) []string {
	if len(types) == 0 {
		return defaultFields[model.UnknownResourceType]
	}
	var names []string
	unique := map[string]bool{}
	for _, rt := range types {
		for _, field := range defaultFields[rt] {
			if unique[field] {
				continue
			}
			unique[field] = true
			names = append(names, field)
		}
	}
	return names
}

func ColumnsByName(names []string) []utils.Column {
	cols := make([]utils.Column, len(names))
	for i, name := range names {
		cols[i] = knownColumnForName(name)
	}
	return cols
}

func knownColumnForName(name string) utils.Column {
	for _, col := range knownColumns {
		if strings.EqualFold(col.Name, name) {
			return col
		}
	}
	return utils.DefaultColumn
}

var knownColumns = []utils.Column{
	{
		Name:      "filename",
		Alignment: utils.Left,
		Width:     -1,
	},
	{
		Name:      "subject",
		Alignment: utils.Left,
		Width:     30,
	},
	{
		Name:      "issuer",
		Alignment: utils.Left,
		Width:     30,
	},
	{
		Name:      "type",
		Alignment: utils.Left,
		Width:     16,
	},
	{
		Name:      "id",
		Alignment: utils.Left,
		Width:     40,
	},
	{
		Name:      "not-after",
		Alignment: utils.Centre,
		Width:     20,
	},
	{
		Name:      "not-before",
		Alignment: utils.Centre,
		Width:     20,
	},
	{
		Name:      "this-update",
		Alignment: utils.Centre,
		Width:     25,
	},
	{
		Name:      "is-ca",
		Alignment: utils.Left,
		Width:     5,
	},
	{
		Name:      "self-signed",
		Alignment: utils.Left,
		Width:     5,
	},
	{
		Name:      "public-key-algorithm",
		Title:     "key algorithm",
		Alignment: utils.Left,
		Width:     11,
	},
	{
		Name:      "version",
		Alignment: utils.Right,
		Width:     8,
	},
	{
		Name:      "number",
		Alignment: utils.Right,
		Width:     8,
	},
}
