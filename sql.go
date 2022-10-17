package md2sql

import (
	"fmt"
	"io"
	"strings"
)

func DumpSQL(w io.Writer, tables []*Table, d Dialect) error {
	_, err := fixRelations(tables, d)
	if err != nil {
		return err
	}

	// table definition
	for i, t := range tables {
		if i != 0 {
			fmt.Fprintf(w, "\n\n")
		}
		fmt.Fprintf(w, "CREATE TABLE %s(\n", t.Name)
		var rows []string
		fkSrcs := make(map[string][]string)
		fkDests := make(map[string][]string)
		var pks []string
		var fkTables []string
		for _, c := range t.Columns {
			if c.PrimaryKey {
				pks = append(pks, c.Name)
				rows = append(rows, fmt.Sprintf("\t%s %s", c.Name, d.PrimaryKeySQLType(c.Type, c.AutoIncrement)))
			} else if c.AssociativeEntity {
				// do nothing
			} else if c.Nullable {
				rows = append(rows, fmt.Sprintf("\t%s %s", c.Name, d.TypeConversion(c.Type)))
			} else {
				rows = append(rows, fmt.Sprintf("\t%s %s NOT NULL", c.Name, d.TypeConversion(c.Type)))
			}
			if c.LinkTable != "" && !c.AssociativeEntity {
				fkSrcs[c.LinkTable] = append(fkSrcs[c.LinkTable], c.Name)
				fkDests[c.LinkTable] = append(fkDests[c.LinkTable], c.LinkColumn)
				fkTables = append(fkTables, c.LinkTable)
			}
		}
		if len(pks) > 0 {
			rows = append(rows, fmt.Sprintf("\tPRIMARY KEY(%s)", strings.Join(pks, ", ")))
		}
		for _, t := range fkTables {
			rows = append(rows, fmt.Sprintf("\tFOREIGN KEY(%s) REFERENCES %s(%s)", strings.Join(fkSrcs[t], ", "), t, strings.Join(fkDests[t], ", ")))
		}
		fmt.Fprintf(w, "%s\n);", strings.Join(rows, ",\n"))

		for _, c := range t.Columns {
			if c.Index {
				fmt.Fprintf(w, "\n\nCREATE UNIQUE INDEX INDEX_%s_%s ON %s(%s);", t.Name, c.Name, t.Name, c.Name)
			}
		}
	}

	// associative entity
	for _, t := range tables {
		var pks []string
		var pkTypes []string
		for _, c := range t.Columns {
			if c.PrimaryKey {
				pks = append(pks, c.Name)
				pkTypes = append(pkTypes, c.Type)
			}
		}

		for _, c := range t.Columns {
			if c.LinkTable != "" && c.AssociativeEntity {
				fmt.Fprintf(w, "\n\nCREATE TABLE %s_%s(\n", t.Name, c.Name)
				var rows []string
				rows = append(rows, fmt.Sprintf("\tid %s PRIMARY KEY", d.PrimaryKeySQLType("", true)))
				var fks []string
				for i, pk := range pks {
					rows = append(rows, fmt.Sprintf("\t%s_%s %s", t.Name, pk, d.PrimaryKeyBaseType(pkTypes[i])))
					fks = append(fks, t.Name+"_"+pk)
				}
				rows = append(rows, fmt.Sprintf("\t%s_%s %s", c.LinkTable, c.LinkColumn, d.PrimaryKeyBaseType(c.Type)))
				rows = append(rows, fmt.Sprintf("\tFOREIGN KEY(%s) REFERENCES %s(%s)", strings.Join(fks, ", "), t.Name, strings.Join(pks, ", ")))
				rows = append(rows, fmt.Sprintf("\tFOREIGN KEY(%s_%s) REFERENCES %s(%s)", c.LinkTable, c.LinkColumn, c.LinkTable, c.LinkColumn))
				fmt.Fprintf(w, "%s\n);", strings.Join(rows, ",\n"))
			}
		}
	}

	return nil
}
