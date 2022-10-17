package md2sql

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("🐙")
}

func TestParseColumn(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    Column
		wantErr bool
	}{
		{
			name: "simple column",
			args: args{
				src: "name: text",
			},
			want: Column{
				Name: "name",
				Type: "text",
			},
		},
		{
			name: "index",
			args: args{
				src: "#email: text",
			},
			want: Column{
				Name:  "email",
				Type:  "text",
				Index: true,
			},
		},
		{
			name: "primary key without type",
			args: args{
				src: "##id",
			},
			want: Column{
				Name:          "id",
				PrimaryKey:    true,
				AutoIncrement: true,
			},
		},
		{
			name: "primary key with type",
			args: args{
				src: "##email: varchar(30)",
			},
			want: Column{
				Name:       "email",
				Type:       "varchar(30)",
				PrimaryKey: true,
			},
		},
		{
			name: "nullable",
			args: args{
				src: "name: text?",
			},
			want: Column{
				Name:     "name",
				Type:     "text",
				Nullable: true,
			},
		},
		{
			name: "foreign key(ok)",
			args: args{
				src: "job: *job.id",
			},
			want: Column{
				Name:       "job",
				LinkTable:  "job",
				LinkColumn: "id",
			},
		},
		{
			name: "foreign key(error)",
			args: args{
				src: "job: *jobid", // no period
			},
			wantErr: true,
		},
		{
			name: "foreign key (nullable)",
			args: args{
				src: "job: *job.id?",
			},
			want: Column{
				Name:       "job",
				LinkTable:  "job",
				LinkColumn: "id",
				Nullable:   true,
			},
		},
		{
			name: "foreign key (associative entity)",
			args: args{
				src: "job: *job.id[]",
			},
			want: Column{
				Name:              "job",
				LinkTable:         "job",
				LinkColumn:        "id",
				AssociativeEntity: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColumn(tt.args.src)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.want, got)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Table
		wantErr bool
	}{
		{
			name: "simple single table",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * name: text
				  * age:  int
				`),
			},
			want: []*Table{
				{
					Name: "User",
					Columns: []*Column{
						{
							Name: "name",
							Type: "text",
						},
						{
							Name: "age",
							Type: "int",
						},
					},
				},
			},
		},
		{
			name: "two tables",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * name: text
				  * age:  int
				* table: Job
				  * name: text
				`),
			},
			want: []*Table{
				{
					Name: "User",
					Columns: []*Column{
						{
							Name: "name",
							Type: "text",
						},
						{
							Name: "age",
							Type: "int",
						},
					},
				},
				{
					Name: "Job",
					Columns: []*Column{
						{
							Name: "name",
							Type: "text",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tt.args.src))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
