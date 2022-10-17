package md2sql

import (
	"fmt"
	"io"
)

var plantumlCN = map[Cardinality][]string{
	ZeroOrOne:  {"|o", "o|"},
	ExactlyOne: {"||", "||"},
	ZeroOrMore: {"}o", "o{"},
}

func DumpPlantUML(w io.Writer, tables []*Table, d Dialect) error {
	relations, err := fixRelations(tables, d)
	if err != nil {
		return err
	}
	io.WriteString(w, "@startuml\n\n")
	for _, t := range tables {
		fmt.Fprintf(w, "entity %s {\n", t.Name)
		for _, c := range t.Columns {
			if c.PrimaryKey {
				if c.AutoIncrement {
					fmt.Fprintf(w, "  *%s:%s <<generated>>\n", c.Name, d.PrimaryKeyBaseType(c.Type))
				} else {
					fmt.Fprintf(w, "  *%s:%s\n", c.Name, d.PrimaryKeyBaseType(c.Type))
				}
			}
		}
		fmt.Fprintf(w, "  --\n")
		for _, c := range t.Columns {
			if c.PrimaryKey {
				continue
			}
			if c.LinkTable != "" {
				if !c.AssociativeEntity {
					if c.Nullable {
						fmt.Fprintf(w, "  %s:%s <<FK>>\n", c.Name, d.PrimaryKeyBaseType(c.Type))
					} else {
						fmt.Fprintf(w, "  *%s:%s <<FK>>\n", c.Name, d.PrimaryKeyBaseType(c.Type))
					}
				}
			} else if c.Nullable {
				fmt.Fprintf(w, "  %s:%s\n", c.Name, d.TypeConversion(c.Type))
			} else {
				fmt.Fprintf(w, "  *%s:%s\n", c.Name, d.TypeConversion(c.Type))
			}
		}
		fmt.Fprintf(w, "}\n\n")
	}
	for _, r := range relations {
		fmt.Fprintf(w, "%s %s--%s %s\n\n", r.FromTable, plantumlCN[r.FromCardinality][0], plantumlCN[r.ToCardinality][1], r.ToTable)
	}
	io.WriteString(w, "@enduml")
	return nil
}
