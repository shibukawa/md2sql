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

func TestDumpGraphviz(t *testing.T) {
	type args struct {
		src  string
		mode ModelType
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
				  * age:  integer
				`),
				mode: PhysicalModel,
			},
			want: TrimIndent(t, `
			digraph erd {
				graph [rankdir=LR, overlap=false, splines=true];
				edge [dir=both];
				node [shape=Mrecord, fontname=verdana, fontsize=9];

				table0 [label=<
					<table border="0" cellspacing="2" cellpadding="0"><tr><td><b>Users</b></td></tr></table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">PK&nbsp;<b>id</b>&nbsp;<i><font color="lightgray">INTEGER</font></i></td></tr>
					</table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">*<b>name</b>&nbsp;<i><font color="lightgray">TEXT</font></i></td></tr>
						<tr><td align="left">*<b>age</b>&nbsp;<i><font color="lightgray">INTEGER</font></i></td></tr>
					</table>>];
			}
			`),
		},
		{
			name: "two table with foreign key",
			args: args{
				src: TrimIndent(t, `
				* table: Users
				  * @id
				  * job: *Jobs.id
				  * name: string
				* table: Jobs
				  * @id
				  * name: string
				`),
			},
			want: TrimIndent(t, `
			digraph erd {
				graph [rankdir=LR, overlap=false, splines=true];
				edge [dir=both];
				node [shape=Mrecord, fontname=verdana, fontsize=9];

				table0 [label=<
					<table border="0" cellspacing="2" cellpadding="0"><tr><td><b>Users</b></td></tr></table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">PK&nbsp;<b>id</b>&nbsp;<i><font color="lightgray">INTEGER</font></i></td></tr>
					</table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">*FK&nbsp;<b>job</b>&nbsp;<i><font color="lightgray">INTEGER</font></i></td></tr>
						<tr><td align="left">*<b>name</b>&nbsp;<i><font color="lightgray">TEXT</font></i></td></tr>
					</table>>];
				table1 [label=<
					<table border="0" cellspacing="2" cellpadding="0"><tr><td><b>Jobs</b></td></tr></table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">PK&nbsp;<b>id</b>&nbsp;<i><font color="lightgray">INTEGER</font></i></td></tr>
					</table>
					|<table border="0" cellspacing="2" cellpadding="0">
						<tr><td align="left">*<b>name</b>&nbsp;<i><font color="lightgray">TEXT</font></i></td></tr>
					</table>>];

				table0 -> table1 [arrowtail=crow, arrowhead=tee];
			}
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
			DumpGraphviz(w, tables, tt.args.mode, PostgreSQL)
			assert.Equal(t, tt.want, w.String())
		})
	}
}
