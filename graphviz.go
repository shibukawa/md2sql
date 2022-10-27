package md2sql

import (
	"fmt"
	"io"
)

type ModelType int

const (
	PhysicalModel ModelType = iota
	LogicalModel
	ConceptualModel
)

var graphvizCN = map[Cardinality]string{
	ZeroOrOne:  "odot",
	ExactlyOne: "tee",
	ZeroOrMore: "crow",
}

func DumpGraphviz(w io.Writer, tables []*Table, m ModelType, d Dialect) error {
	relations, err := fixRelations(tables, d)
	if err != nil {
		return err
	}

	tableIDs := make(map[string]string)

	io.WriteString(w, trimIndent(`
		digraph erd {
			graph [rankdir=LR, overlap=false, splines=true];
			edge [dir=both];
			node [shape=Mrecord, fontname=verdana, fontsize=9];
		
		`, ""))
	for i, t := range tables {
		tableIDs[t.Name] = fmt.Sprintf("table%d", i)
		fmt.Fprintf(w, trimIndent(`
			table%d [label=<
				<table border="0" cellspacing="2" cellpadding="0"><tr><td><b>%s</b></td></tr></table>`, "\t"), i, t.Name)

		// primarykeys
		fmt.Fprintf(w, "\n\t\t|<table border=\"0\" cellspacing=\"2\" cellpadding=\"0\">\n")
		for _, c := range t.Columns {
			if !c.PrimaryKey {
				continue
			}
			fmt.Fprintf(w, "\t\t\t<tr><td align=\"left\">PK&nbsp;<b>%s</b>&nbsp;<i><font color=\"lightgray\">%s</font></i></td></tr>\n", c.Name, d.PrimaryKeyBaseType(c.Type))
		}
		fmt.Fprintf(w, "\t\t</table>\n")

		// other fields
		fmt.Fprintf(w, "\t\t|<table border=\"0\" cellspacing=\"2\" cellpadding=\"0\">\n")
		for _, c := range t.Columns {
			if c.PrimaryKey {
				continue
			}
			tn := ""
			cst := ""
			if c.LinkTable != "" {
				if !c.AssociativeEntity {
					tn = d.PrimaryKeyBaseType(c.Type)
					if c.Nullable {
						cst = "FK&nbsp;"
					} else {
						cst = "*FK&nbsp;"
					}
				}
			} else if c.Nullable {
				tn = d.TypeConversion(c.Type)
			} else {
				cst = "*"
				tn = d.TypeConversion(c.Type)
			}
			fmt.Fprintf(w, "\t\t\t<tr><td align=\"left\">%s<b>%s</b>&nbsp;<i><font color=\"lightgray\">%s</font></i></td></tr>\n", cst, c.Name, tn)
		}
		fmt.Fprintf(w, "\t\t</table>>];\n")
	}
	for _, r := range relations {
		fmt.Fprintf(w, "\n\t%s -> %s [arrowtail=%s, arrowhead=%s];\n", tableIDs[r.FromTable], tableIDs[r.ToTable], graphvizCN[r.FromCardinality], graphvizCN[r.ToCardinality])
	}
	io.WriteString(w, "}")
	return nil
}
