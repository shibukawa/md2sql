package md2sql

import (
	"strings"
)

//go:generate go run github.com/dmarkham/enumer -type=Dialect

type Dialect int

const (
	PostgreSQL Dialect = iota
	MySQL
	SQLite
)

func ToDialect(src string) Dialect {
	switch strings.ToLower(src) {
	case "postgres":
		return PostgreSQL
	case "postgresql":
		return PostgreSQL
	case "pg":
		return PostgreSQL
	case "mysql":
		return MySQL
	case "maria":
		return MySQL
	case "mariadb":
		return MySQL
	case "sqlite":
		return SQLite
	}
	return PostgreSQL
}

func (d Dialect) PrimaryKeySQLType(t string, autoIncrement bool) string {
	if t == "" {
		if autoIncrement {
			switch d {
			case PostgreSQL:
				return "SERIAL"
			case MySQL:
				return "SERIAL"
			case SQLite:
				return "INTEGER AUTOINCREMENT"
			}
		} else {
			return d.PrimaryKeyBaseType("")
		}
	}
	return d.TypeConversion(t)
}

func (d Dialect) EnableForeignKey(hasForeignKey bool) string {
	if d == SQLite && hasForeignKey {
		return "PRAGMA foreign_keys = ON;\n\n"
	}
	return ""
}

func (d Dialect) PrimaryKeyBaseType(t string) string {
	if t == "" {
		switch d {
		case PostgreSQL:
			return "INTEGER"
		case MySQL:
			return "INTEGER"
		case SQLite:
			return "INTEGER"
		}
	}
	return d.TypeConversion(t)
}

func (d Dialect) TypeConversion(t string) string {
	switch strings.ToLower(t) {
	case "string":
		fallthrough
	case "text":
		return "TEXT"
	case "blob":
		fallthrough
	case "lob":
		switch d {
		case PostgreSQL:
			return "BYTEA"
		case MySQL:
			return "BLOB"
		case SQLite:
			return "BLOB"
		}
	}
	return strings.ToUpper(t)
}
