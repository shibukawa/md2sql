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
				* table: User
				  * @id
				  * name: text
				  * age:  int?
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity User {
			  *id:INTEGER <<generated>>
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
				* table: User
				  * @id
				  * job: *Job.id
				* table: Job
				  * @id
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity User {
			  *id:INTEGER <<generated>>
			  --
			  *job:INTEGER <<FK>>
			}

			entity Job {
			  *id:INTEGER <<generated>>
			  --
			}

			User }o--|| Job

			@enduml
			`),
		},
		{
			name: "two table with associative entity",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * @id
				  * job: *Job.id[]
				* table: Job
				  * @id
				`),
			},
			want: TrimIndent(t, `
			@startuml

			entity User {
			  *id:INTEGER <<generated>>
			  --
			}

			entity Job {
			  *id:INTEGER <<generated>>
			  --
			}

			User }o--o{ Job

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
			assert.Equal(t, tt.want, w.String())
		})
	}
}
