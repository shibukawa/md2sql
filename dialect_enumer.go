// Code generated by "enumer -type=Dialect"; DO NOT EDIT.

package md2sql

import (
	"fmt"
	"strings"
)

const _DialectName = "PostgreSQLMySQLSQLite"

var _DialectIndex = [...]uint8{0, 10, 15, 21}

const _DialectLowerName = "postgresqlmysqlsqlite"

func (i Dialect) String() string {
	if i < 0 || i >= Dialect(len(_DialectIndex)-1) {
		return fmt.Sprintf("Dialect(%d)", i)
	}
	return _DialectName[_DialectIndex[i]:_DialectIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DialectNoOp() {
	var x [1]struct{}
	_ = x[PostgreSQL-(0)]
	_ = x[MySQL-(1)]
	_ = x[SQLite-(2)]
}

var _DialectValues = []Dialect{PostgreSQL, MySQL, SQLite}

var _DialectNameToValueMap = map[string]Dialect{
	_DialectName[0:10]:       PostgreSQL,
	_DialectLowerName[0:10]:  PostgreSQL,
	_DialectName[10:15]:      MySQL,
	_DialectLowerName[10:15]: MySQL,
	_DialectName[15:21]:      SQLite,
	_DialectLowerName[15:21]: SQLite,
}

var _DialectNames = []string{
	_DialectName[0:10],
	_DialectName[10:15],
	_DialectName[15:21],
}

// DialectString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DialectString(s string) (Dialect, error) {
	if val, ok := _DialectNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DialectNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Dialect values", s)
}

// DialectValues returns all values of the enum
func DialectValues() []Dialect {
	return _DialectValues
}

// DialectStrings returns a slice of all String values of the enum
func DialectStrings() []string {
	strs := make([]string, len(_DialectNames))
	copy(strs, _DialectNames)
	return strs
}

// IsADialect returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Dialect) IsADialect() bool {
	for _, v := range _DialectValues {
		if i == v {
			return true
		}
	}
	return false
}
