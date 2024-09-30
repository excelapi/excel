package xlsx

import (
	"fmt"
	"math"
	"strconv"
)

type ColumnPrediction struct {
	FileIdx    int
	ColumnName string
	ColumnType string
	Decorators []string
	NotNull    bool
	PrimaryKey bool
}

type TablePrediction struct {
	TableName string
	Columns   []ColumnPrediction
}

func (ws *Worksheet) PredictTable() (*TablePrediction, error) {
	headerIdx, err := ws.headerRowIndex()
	if err != nil {
		return nil, err
	}

	header, err := ws.suggestHeader(headerIdx)
	if err != nil {
		return nil, err
	}

	dataRowIdx := headerIdx + 1
	columnPredictions := ws.predictColumnTypes(dataRowIdx, header)
	fmt.Println(columnPredictions)

	// predictedPK := columnNames[0]

	return &TablePrediction{}, nil
}

// index of the header
func (ws *Worksheet) headerRowIndex() (int, error) {
	data := ws.Sheet.Rows
	if len(data) == 0 {
		return 0, fmt.Errorf("no rows found in sheet")
	}

	fmt.Println("Number of rows found:", len(data))

	// find first row with cells
	for i := 0; i < len(data); i++ {
		row := data[i]
		// fmt.Printf("Row %d has %d cells\n", i, len(row.Cells))

		// if row contains cells
		if len(row.Cells) > 0 {
			return i, nil // return the index
		}
	}

	return 0, fmt.Errorf("unable to locate potential header row")
}

// headers based on row index
// ? is this working?
func (ws *Worksheet) suggestHeader(idx int) ([]string, error) {
	header := ws.Sheet.Rows[idx]
	headerStr := []string{}

	for _, cell := range header.Cells {
		if cell.Type == "s" {
			idx, _ := strconv.Atoi(cell.Value)
			headerStr = append(headerStr, ws.SS.StringItems[idx].Text)
		} else {
			return nil, fmt.Errorf("all header names must be strings")
		}
	}
	return headerStr, nil
}

// predict the types based
func (ws *Worksheet) predictColumnTypes(dataStart int, header []string) []ColumnPrediction {

	columnPredictions := []ColumnPrediction{}

	// i loops through each column once
	// datastart is the row that we think data starts on
	for i := 0; i < len(header); i++ {
		// starting cell
		startingCell := ws.Sheet.Rows[dataStart].Cells[i]

		// loop through each row in a single column to get the cp (column prediction)
		cp := ws.determineColumnType(dataStart, i, "", startingCell, &map[string]bool{}, &ColumnPrediction{FileIdx: i, NotNull: true, PrimaryKey: true})

		// add column name to the cp
		cp.ColumnName = header[i]

		columnPredictions = append(columnPredictions, *cp)
	}

	return columnPredictions
}

func (c *Cell) dataType() string {
	if c.Type == "s" {
		return "string"
	} else if c.Type == "b" {
		return "boolean"
	} else {
		d, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			fmt.Println("datatype error: ", err.Error())
			return "string"
		}

		if math.Mod(d, 1) == 0 {
			return "integer"
		}
		return "float64"
	}
}

// ? determine precision !! BIG DEAL
// ? determine uniqueness
// ? determine empty values
// ? is nullable?
// ? is primary key eligible
func (ws *Worksheet) determineColumnType(row, col int, prevType string, cell Cell, dups *map[string]bool, cp *ColumnPrediction) *ColumnPrediction {

	// check if pk and NotNull are still options
	if cp.PrimaryKey || cp.NotNull {
		// get value
		value := ws.getValue(&cell)

		if (*dups)[value] {
			cp.PrimaryKey = false
		} else {
			(*dups)[value] = true
		}

		// check if value is empty
		// if empty, remove pk and not null eligable
		if value == "" {
			// remove pk and not null
			cp.PrimaryKey = false
			cp.NotNull = false
		}
	}
	// determine cell type
	cellType := cell.dataType()

	// if precision is bigger than before, replace
	// else just use old precision

	// if it changes, we're likely to just change it to a string
	if prevType != cellType {
		// we want to allow the int -> float change
		if prevType != "integer" && cellType != "float64" {
			// new type is just going to be a string
			// check that we can go another
			if row+1 < len(ws.Sheet.Rows) {
				return ws.determineColumnType(row+1, col, "string", ws.Sheet.Rows[row+1].Cells[col], dups, cp)
			} else {
				cp.ColumnType = dateTypeToColumnType(cellType)
				return cp
			}
		}
	}

	if row+1 < len(ws.Sheet.Rows) {
		// type didn't change
		return ws.determineColumnType(row+1, col, cellType, ws.Sheet.Rows[row+1].Cells[col], dups, cp)
	} else {
		cp.ColumnType = dateTypeToColumnType(cellType)
		return cp
	}
}

func dateTypeToColumnType(t string) string {
	switch t {
	case "string":
		return "TEXT"
	case "integer":
		return "INTEGER"
	case "float64":
		return "DOUBLE PRECISION"
	default:
		return "TEXT"
	}
}
