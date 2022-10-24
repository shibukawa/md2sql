//go:build wasm

package main

import (
	"bytes"
	"strings"
	"syscall/js"

	"github.com/shibukawa/md2sql"
)

func ConvertToSQL(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{
			"ok":      false,
			"message": "first argument should be markdown source.",
		}
	}
	dialect := "postgres"
	if len(args) == 2 {
		dialect = args[1].String()
	}
	tables, err := md2sql.Parse(strings.NewReader(args[0].String()))
	if err != nil {
		return map[string]any{
			"ok":      false,
			"message": err.Error(),
		}
	}
	var buf bytes.Buffer
	md2sql.DumpSQL(&buf, tables, md2sql.ToDialect(dialect))
	return map[string]any{
		"ok":     true,
		"result": buf.String(),
	}
}

func ConvertToMermaid(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{
			"ok":      false,
			"message": "first argument should be markdown source.",
		}
	}
	tables, err := md2sql.Parse(strings.NewReader(args[0].String()))
	if err != nil {
		return map[string]any{
			"ok":      false,
			"message": err.Error(),
		}
	}
	var buf bytes.Buffer
	md2sql.DumpMermaid(&buf, tables, md2sql.PostgreSQL)
	return map[string]any{
		"ok":     true,
		"result": buf.String(),
	}
}

func ConvertToPlantUML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{
			"ok":      false,
			"message": "first argument should be markdown source.",
		}
	}
	tables, err := md2sql.Parse(strings.NewReader(args[0].String()))
	if err != nil {
		return map[string]any{
			"ok":      false,
			"message": err.Error(),
		}
	}
	var buf bytes.Buffer
	md2sql.DumpPlantUML(&buf, tables, md2sql.PostgreSQL)
	return map[string]any{
		"ok":     true,
		"result": buf.String(),
	}
}

func main() {
	c := make(chan struct{})
	js.Global().Set("md2sql", js.ValueOf(map[string]any{
		"toSQL":      js.FuncOf(ConvertToSQL),
		"toMermaid":  js.FuncOf(ConvertToMermaid),
		"toPlantUML": js.FuncOf(ConvertToPlantUML),
	}))
	<-c
}
