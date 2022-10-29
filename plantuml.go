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

var plantTheme = `
!define ENTITY #A5D6A7-E0E0E0 
!define MASTER_ENTITY #64B5F6-BDBDBD
!define TRANSACTION_ENTITY #FFB74D-BDBDBD
!define SUMMARY_ENTITY #FFF176-BDBDBD
!define WORK_ENTITY #EF9A9A-BDBDBD
!define VIEW_ENTITY #CE93D8-BDBDBD

!define D_ENTITY #E8F5E9-FAFAFA
!define D_MASTER_ENTITY #E3F2FD-E0E0E0
!define D_TRANSACTION_ENTITY #FFF3E0-E0E0E0
!define D_SUMMARY_ENTITY #FFFDE7-E0E0E0
!define D_WORK_ENTITY #FFEBEE-E0E0E0
!define D_VIEW_ENTITY #F3E5F5-E0E0E0
!define D_ASSOCIATIVE_ENTITY #BCAAA4-BDBDBD

!define ENTITY_MARK_COLOR 66BB6A
!define MASTER_MARK_COLOR 42A5F5
!define TRANSACTION_MARK_COLOR FFA726
!define SUMMARY_MARK_COLOR FFEE58
!define WORK_MARK_COLOR EF5350
!define VIEW_MARK_COLOR AB47BC
!define ASSOCIATIVE_MARK_COLOR A1887F

skinparam class {
    BackgroundColor ENTITY
    BorderColor Black
    ArrowColor Black
}
`

/*
entity "ユーザー" as User <<A,ASSOCIATIVE_MARK_COLOR>> D_ASSOCIATIVE_ENTITY {
  +id:INTEGER <<PK>> <<generated>>
  --
  *name:TEXT
  *email:TEXT
  *age:INTEGER
  #jobs:INTEGER <<FK>>
}
*/

var theme = map[TableType]map[bool]string{
	EntityTable: {
		true:  "<<E,ENTITY_MARK_COLOR>> ENTITY",
		false: "<<E,ENTITY_MARK_COLOR>> D_ENTITY",
	},
	MasterTable: {
		true:  "<<M,MASTER_MARK_COLOR>> MASTER_ENTITY",
		false: "<<M,MASTER_MARK_COLOR>> D_MASTER_ENTITY",
	},
	TransactionTable: {
		true:  "<<T,TRANSACTION_MARK_COLOR>> TRANSACTION_ENTITY",
		false: "<<T,TRANSACTION_MARK_COLOR>> D_TRANSACTION_ENTITY",
	},
	SummaryTable: {
		true:  "<<S,SUMMARY_MARK_COLOR>> SUMMARY_ENTITY",
		false: "<<S,SUMMARY_MARK_COLOR>> D_SUMMARY_ENTITY",
	},
	WorkTable: {
		true:  "<<W,WORK_MARK_COLOR>> VIEW_ENTITY",
		false: "<<W,WORK_MARK_COLOR>> D_VIEW_ENTITY",
	},
	View: {
		true:  "<<V,VIEW_MARK_COLOR>> ASSOCIATIVE_ENTITY",
		false: "<<V,VIEW_MARK_COLOR>> D_ASSOCIATIVE_ENTITY",
	},
	AssociativeEntity: {
		true:  "<<A,ASSOCIATIVE_MARK_COLOR>> D_ASSOCIATIVE_ENTITY",
		false: "<<A,ASSOCIATIVE_MARK_COLOR>> D_ASSOCIATIVE_ENTITY",
	},
}

func DumpPlantUML(w io.Writer, tables []*Table, d Dialect) error {
	relations, err := fixRelations(tables, d)
	if err != nil {
		return err
	}

	tableIDs := make(map[string]string)

	fmt.Fprintf(w, "@startuml\n\n%s\n", plantTheme)
	for i, t := range tables {
		tableIDs[t.Name] = fmt.Sprintf("table%d", i)
		fmt.Fprintf(w, "entity table%d as \"%s\" %s {\n", i, t.Name, theme[t.Type][t.Independent])
		for _, c := range t.Columns {
			if c.PrimaryKey {
				if c.AutoIncrement {
					fmt.Fprintf(w, "  *%s:%s <<PK>>\n", c.Name, d.PrimaryKeyBaseType(c.Type))
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
		fmt.Fprintf(w, "%s %s--%s %s\n\n", tableIDs[r.FromTable], plantumlCN[r.FromCardinality][0], plantumlCN[r.ToCardinality][1], tableIDs[r.ToTable])
	}
	io.WriteString(w, "@enduml")
	return nil
}
