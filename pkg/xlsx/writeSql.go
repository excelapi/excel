package xlsx

import (
	"fmt"
	"strings"
)

type ColumnPrediction struct {
	FileIdx    int
	ColumnName string
	ColumnType string
	Decorators []string
	PrimaryKey bool
}

type TablePrediction struct {
	TableName string
	Columns   []ColumnPrediction
}

func (ws *Worksheet) WriteSQL(tp *TablePrediction) (string, error) {
	// Build CREATE statement

	createStmt, err := createTable(tp)
	if err != nil {
		return "", err
	}

	fmt.Println(createStmt)

	return "s3/relative/path.file", nil
}

func createTable(tp *TablePrediction) (string, error) {
	// todo: null checks
	// todo: weird cases
	// todo: check that all slices are the same len()

	// head
	sql := fmt.Sprintf("CREATE TABLE %v (\n", tp.TableName)
	pk := "PRIMARY KEY ("

	for _, col := range tp.Columns {
		// column name
		sql += col.ColumnName + " "

		// column type
		sql += col.ColumnType + " "

		// decorators
		if len(col.Decorators) > 0 {
			for _, d := range col.Decorators {
				sql += d + " "
			}
		}

		// remove " " space and add comma ,
		sql = strings.TrimSuffix(sql, " ") + ",\n"

		// primary key?
		if col.PrimaryKey {
			pk += col.ColumnName + ", "
		}
	}

	// todo: check that we actually have a pk
	// clean primary key
	pk = strings.TrimSuffix(pk, ", ") + ")\n"

	// add Primary key
	sql += pk + ")"

	return sql, nil
}

func ParseTablePrediction(body string) (*TablePrediction, error) {

	// todo
	return &TablePrediction{}, nil
}
