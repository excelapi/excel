package xlsx

import (
	"fmt"
	"os"
	"strings"
)

func (ws *Worksheet) WriteSQL(tp *TablePrediction) (string, error) {
	// Build CREATE statement

	createStmt, err := createTable(tp)
	if err != nil {
		return "", err
	}
	// fmt.Println(createStmt)

	insertStmt, err := createInsert(tp, ws)
	if err != nil {
		return "", err
	}
	// fmt.Println(insertStmt)

	var sql string = createStmt + "\n" + insertStmt

	// ! for testing
	err = os.WriteFile("../sql/Test.sql", []byte(sql), 0644)
	if err != nil {
		panic("error writing query to file: " + err.Error())
	}

	return "s3/relative/path.file", nil
}

func createInsert(tp *TablePrediction, ws *Worksheet) (string, error) {

	sql := fmt.Sprintf("INSERT INTO\n\t%v", tp.TableName)
	sql += " ("

	// sort tp columns by FileIdx
	// sort.Sort(tp.Columns)
	// bubble sort
	// or quick sort
	for _, c := range tp.Columns {
		sql += c.ColumnName + ", "
	}
	sql = strings.TrimSuffix(sql, ", ") + ")"
	sql += "\nVALUES\n"

	for i := 1; i < len(ws.Sheet.Rows); i++ {
		row := ws.Sheet.Rows[i]
		sqlRow := "("

		for _, col := range tp.Columns {
			// find the cell for this column
			cell := row.Cells[col.FileIdx]
			if strings.Contains(col.ColumnType, "varchar") {
				sqlRow += "'" + ws.getValue(&cell) + "', "
			} else {
				sqlRow += ws.getValue(&cell) + ", "
			}
		}
		// clean row
		sqlRow = strings.TrimSuffix(sqlRow, ", ") + "),\n"
		sql += sqlRow
	}
	// clean up sql
	sql = strings.TrimSuffix(sql, ",\n") + ";"

	return sql, nil
}

func createTable(tp *TablePrediction) (string, error) {
	// todo: null checks
	// todo: weird cases
	// todo: check that all slices are the same len()

	// head
	sql := fmt.Sprintf("CREATE TABLE %v (\n", tp.TableName)
	pk := "\tPRIMARY KEY ("

	for _, col := range tp.Columns {
		// tab
		sql += "\t"

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
	sql += pk + ");"

	return sql, nil
}
