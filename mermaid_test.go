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

func TestDumpMermaid(t *testing.T) {
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
				  * ##id
				  * name: text
				  * age:  int
				`),
			},
			want: TrimIndent(t, `
			erDiagram

			User {
			  INTEGER id PK
			  TEXT name
			  INT age
			}
			`),
		},
		{
			name: "two table with foreign key",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * ##id
				  * job: *Job.id
				* table: Job
				  * ##id
				`),
			},
			want: TrimIndent(t, `
			erDiagram

			User {
			  INTEGER id PK
			  INTEGER job FK
			}

			Job {
			  INTEGER id PK
			}

			User }o--|| Job : job
			`),
		},
		{
			name: "two table with foreign key (string)",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * ##id
				  * job: *Job.id
				* table: Job
				  * ##id: string
				`),
			},
			want: TrimIndent(t, `
			erDiagram

			User {
			  INTEGER id PK
			  TEXT job FK
			}

			Job {
			  TEXT id PK
			}

			User }o--|| Job : job
			`),
		},
		{
			name: "two table with associative entity",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * ##id
				  * job: *Job.id[]
				* table: Job
				  * ##id: string
				`),
			},
			want: TrimIndent(t, `
			erDiagram

			User {
			  INTEGER id PK
			}

			Job {
			  TEXT id PK
			}

			User }o--o{ Job : job
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
			DumpMermaid(w, tables, PostgreSQL)
			assert.Equal(t, tt.want, w.String())
		})
	}
}
