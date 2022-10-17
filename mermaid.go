package md2sql

import (
	"fmt"
	"io"
)

var mermaidCN = map[Cardinality][]string{
	ZeroOrOne:  {"|o", "o|"},
	ExactlyOne: {"||", "||"},
	ZeroOrMore: {"}o", "o{"},
}

func DumpMermaid(w io.Writer, tables []*Table, d Dialect) error {
	relations, err := fixRelations(tables, d)
	if err != nil {
		return err
	}
	io.WriteString(w, "erDiagram\n\n")
	for i, t := range tables {
		if i != 0 {
			io.WriteString(w, "\n\n")
		}
		fmt.Fprintf(w, "%s {\n", t.Name)
		for _, c := range t.Columns {
			if c.PrimaryKey {
				fmt.Fprintf(w, "  %s %s PK\n", d.PrimaryKeyBaseType(c.Type), c.Name)
			} else if c.LinkTable != "" {
				if !c.AssociativeEntity {
					if c.Nullable {
						fmt.Fprintf(w, "  %s? %s FK\n", d.PrimaryKeyBaseType(c.Type), c.Name)
					} else {
						fmt.Fprintf(w, "  %s %s FK\n", d.PrimaryKeyBaseType(c.Type), c.Name)
					}
				}
			} else if c.Nullable {
				fmt.Fprintf(w, "  %s? %s\n", d.TypeConversion(c.Type), c.Name)
			} else {
				fmt.Fprintf(w, "  %s %s\n", d.TypeConversion(c.Type), c.Name)
			}
		}
		io.WriteString(w, "}")
	}
	for _, r := range relations {
		fmt.Fprintf(w, "\n\n%s %s--%s %s : %s", r.FromTable, mermaidCN[r.FromCardinality][0], mermaidCN[r.ToCardinality][1], r.ToTable, r.Label)
	}
	return nil
}
