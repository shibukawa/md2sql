package md2sql

import (
	"fmt"
	"io"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type TableType int

const (
	EntityTable TableType = iota
	MasterTable
	TransactionTable
	WorkTable
	SummaryTable
	View
	AssociativeEntity
)

var label2tableType = map[string]TableType{
	"table":             EntityTable,
	"master":            MasterTable,
	"mastertable":       MasterTable,
	"tran":              TransactionTable,
	"transaction":       TransactionTable,
	"transactiontable":  TransactionTable,
	"work":              WorkTable,
	"summary":           SummaryTable,
	"view":              View,
	"associativeentity": AssociativeEntity,
}

type Table struct {
	Type        TableType
	Independent bool
	Name        string
	Columns     []*Column
}

type Column struct {
	Name                 string
	Type                 string
	LinkTable            string
	LinkColumn           string
	PrimaryKey           bool
	AutoIncrement        bool
	Index                bool
	Nullable             bool
	AssociativeEntity    bool
	ForeignKeyConstraint bool
}

func ParseColumn(src string) (*Column, error) {
	var result Column
	before, after, ok := strings.Cut(src, ":")
	if ok {
		before = strings.TrimSpace(before)
		after = strings.TrimSpace(after)
		if strings.HasPrefix(before, "@") {
			before = strings.TrimPrefix(before, "@")
			result.PrimaryKey = true
		} else if strings.HasPrefix(before, "$") {
			before = strings.TrimPrefix(before, "$")
			result.Index = true
		}
		if strings.HasPrefix(after, "*") {
			after = strings.TrimPrefix(after, "*")
			if strings.HasSuffix(after, "?") {
				after = strings.TrimSuffix(after, "?")
				result.Nullable = true
			} else if strings.HasSuffix(after, "[]") {
				after = strings.TrimSuffix(after, "[]")
				result.AssociativeEntity = true
			}
			table, column, ok := strings.Cut(after, ".")
			if !ok {
				return nil, fmt.Errorf("foreign key definition should be 'table.column', but no period found: %s", after)
			}
			result.LinkTable = table
			result.LinkColumn = column
			after = ""
		}
		if strings.HasSuffix(after, "?") {
			after = strings.TrimSuffix(after, "?")
			result.Nullable = true
		}
		result.Name = before
		result.Type = after
	} else {
		if strings.HasPrefix(before, "@") {
			before = strings.TrimPrefix(before, "@")
			result.PrimaryKey = true
			result.AutoIncrement = true
			result.Name = before
		}
	}
	return &result, nil
}

func Parse(r io.Reader) ([]*Table, error) {
	markdown := goldmark.New()
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	reader := text.NewReader(b)
	n := markdown.Parser().Parse(reader)
	var tables []*Table
	err = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() == ast.KindListItem && entering {
			if n.ChildCount() == 2 && n.LastChild().Kind() == ast.KindList { // nested
				label := string(n.FirstChild().Text(b))
				if t, name, ok := strings.Cut(label, ":"); ok {
					independent := true
					if strings.HasPrefix(t, "_") || strings.HasPrefix(t, "-") {
						t = strings.TrimLeft(t, "_-")
						independent = false
					}
					if tt, ok := label2tableType[strings.ToLower(t)]; ok {
						var columns []*Column
						c := n.LastChild().FirstChild()
						for c != nil {
							column, err := ParseColumn(string(c.FirstChild().Text(b)))
							if err != nil {
								return ast.WalkStop, err
							}
							columns = append(columns, column)
							c = c.NextSibling()
						}
						tables = append(tables, &Table{
							Type:        tt,
							Independent: independent,
							Name:        strings.TrimSpace(name),
							Columns:     columns,
						})
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	return tables, nil
}

type Cardinality int

const (
	ZeroOrOne Cardinality = iota
	ExactlyOne
	ZeroOrMore
)

type Relation struct {
	FromTable       string
	FromCardinality Cardinality
	ToTable         string
	ToCardinality   Cardinality
	Label           string
}

func fixRelations(tables []*Table, d Dialect) ([]*Relation, error) {
	tmap := make(map[string]*Table)
	cmap := make(map[string]*Column)

	key := func(table, column string) string {
		return table + "/**/" + column
	}
	for _, t := range tables {
		tmap[t.Name] = t
		for _, c := range t.Columns {
			cmap[key(t.Name, c.Name)] = c
		}
	}

	var result []*Relation
	// fill type
	for _, t := range tables {
		for _, c := range t.Columns {
			if c.LinkTable != "" {
				rel := &Relation{
					FromTable:       t.Name,
					FromCardinality: ZeroOrMore,
					ToTable:         c.LinkTable,
					Label:           c.Name,
				}
				if tc, ok := cmap[key(c.LinkTable, c.LinkColumn)]; ok {
					c.Type = d.PrimaryKeyBaseType(tc.Type)
				} else {
					c.Type = "INTEGER" // fill dummy
				}
				if c.AssociativeEntity {
					rel.ToCardinality = ZeroOrMore
					rel.FromCardinality = ZeroOrMore
				} else {
					if c.Nullable {
						rel.ToCardinality = ZeroOrOne
					} else {
						rel.ToCardinality = ExactlyOne
					}
				}
				result = append(result, rel)
			}
		}
	}

	return result, nil
}
