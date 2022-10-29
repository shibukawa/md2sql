package md2sql

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("🐙")
}

func TestDumpPlantUML(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "single table",
			args: args{
				src: TrimIndent(t, `
				* table: Users
				  * @id
				  * name: text
				  * age:  int?
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity table0 as "Users" <<E,ENTITY_MARK_COLOR>> ENTITY {
			  *id:INTEGER <<PK>>
			  --
			  *name:TEXT
			  age:INT
			}

			@enduml`),
		},
		{
			name: "two table with foreign key",
			args: args{
				src: TrimIndent(t, `
				* table: Users
				  * @id
				  * job: *Jobs.id
				* table: Jobs
				  * @id
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity table0 as "Users" <<E,ENTITY_MARK_COLOR>> ENTITY {
			  *id:INTEGER <<PK>>
			  --
			  *job:INTEGER <<FK>>
			}

			entity table1 as "Jobs" <<E,ENTITY_MARK_COLOR>> ENTITY {
			  *id:INTEGER <<PK>>
			  --
			}

			table0 }o--|| table1

			@enduml
			`),
		},
		{
			name: "two table with associative entity",
			args: args{
				src: TrimIndent(t, `
				* table: Users
				  * @id
				  * job: *Jobs.id[]
				* table: Jobs
				  * @id
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity table0 as "Users" <<E,ENTITY_MARK_COLOR>> ENTITY {
			  *id:INTEGER <<PK>>
			  --
			}

			entity table1 as "Jobs" <<E,ENTITY_MARK_COLOR>> ENTITY {
			  *id:INTEGER <<PK>>
			  --
			}

			table0 }o--o{ table1

			@enduml
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tables, err := Parse(strings.NewReader(tt.args.src))
			assert.NoError(t, err)
			if err != nil {
				return
			}
			DumpPlantUML(w, tables, PostgreSQL)
			src := w.String()
			assert.True(t, strings.Contains(src, plantTheme))
			src = strings.ReplaceAll(src, plantTheme+"\n", "")
			assert.Equal(t, tt.want, src)
		})
	}
}
