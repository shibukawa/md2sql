package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/shibukawa/md2sql"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("🐙 ")
}

var (
	dialect = kingpin.Flag("dialect", "SQL dialect").Short('d').Default("postgres").Enum("postgres", "mysql", "sqlite")
	format  = kingpin.Flag("format", "Output format").Short('f').Default("sql").Enum("sql", "mermaid", "plantuml")
	output  = kingpin.Flag("output", "Output file").Short('o').File()
	source  = kingpin.Arg("src", "source file").ExistingFile()
)

var dummy = `
# title

* table: Person
	* ##id
	* name: string
	* age: integer

* table: work
	* ##id
`

func main() {
	kingpin.Parse()

	if *output == nil {
		output = &os.Stdout
	} else {
		defer (*output).Close()
	}

	var src io.Reader
	if source == nil {
		src = os.Stdin
	} else {
		sf, err := os.Open(*source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "file open error: %s", err.Error())
			os.Exit(1)
		}
		defer sf.Close()
		src = sf
	}

	tables, err := md2sql.Parse(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %s", err.Error())
		os.Exit(1)
	}
	d := md2sql.ToDialect(*dialect)
	switch *format {
	case "sql":
		md2sql.DumpSQL(*output, tables, d)
	case "mermaid":
		md2sql.DumpMermaid(*output, tables, d)
	case "plantuml":
		md2sql.DumpPlantUML(*output, tables, d)
	}
}
