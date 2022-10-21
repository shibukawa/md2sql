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

func TestSQL(t *testing.T) {
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
				  * name: string
				  * age:  integer?
				`),
			},
			want: TrimIndent(t, `
			CREATE TABLE User(
				id SERIAL,
				name TEXT NOT NULL,
				age INTEGER,
				PRIMARY KEY(id)
			);`),
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
			CREATE TABLE User(
				id SERIAL,
				job INTEGER NOT NULL,
				PRIMARY KEY(id),
				FOREIGN KEY(job) REFERENCES Job(id)
			);

			CREATE TABLE Job(
				id SERIAL,
				PRIMARY KEY(id)
			);
			`),
		},
		{
			name: "two table with primary foreign key",
			args: args{
				src: TrimIndent(t, `
				* table: Users
				  * @id
				  * name: string
				* table: Tasks
				  * @id: *Users.id
				`),
			},
			want: TrimIndent(t, `
			CREATE TABLE Users(
				id SERIAL,
				name TEXT NOT NULL,
				PRIMARY KEY(id)
			);

			CREATE TABLE Tasks(
				id INTEGER,
				PRIMARY KEY(id),
				FOREIGN KEY(id) REFERENCES Users(id)
			);
			`),
		},
		{
			name: "table with primary keys",
			args: args{
				src: TrimIndent(t, `
				* table: Logs
				  * @task:  integer
				  * @index: integer
				  * log: string
				`),
			},
			want: TrimIndent(t, `
			CREATE TABLE Logs(
				task INTEGER,
				index INTEGER,
				log TEXT NOT NULL,
				PRIMARY KEY(task, index)
			);
			`),
		},
		{
			name: "two table with associative entity",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * @id
				  * jobs: *Job.id[]
				* table: Job
				  * @id
				`),
			},
			want: TrimIndent(t, `
			CREATE TABLE User(
				id SERIAL,
				PRIMARY KEY(id)
			);

			CREATE TABLE Job(
				id SERIAL,
				PRIMARY KEY(id)
			);

			CREATE TABLE User_jobs(
				id SERIAL PRIMARY KEY,
				User_id INTEGER,
				Job_id INTEGER,
				FOREIGN KEY(User_id) REFERENCES User(id),
				FOREIGN KEY(Job_id) REFERENCES Job(id)
			);
			`),
		},
		{
			name: "index",
			args: args{
				src: TrimIndent(t, `
				* table: User
				  * @id
				  * name: string
				  * $email: string
				  * age:  integer?
				`),
			},
			want: TrimIndent(t, `
			CREATE TABLE User(
				id SERIAL,
				name TEXT NOT NULL,
				email TEXT NOT NULL,
				age INTEGER,
				PRIMARY KEY(id)
			);
			
			CREATE UNIQUE INDEX INDEX_User_email ON User(email);
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
			DumpSQL(w, tables, PostgreSQL)
			assert.Equal(t, tt.want, w.String())
		})
	}
}
