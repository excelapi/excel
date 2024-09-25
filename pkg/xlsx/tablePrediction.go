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
	columnTypes := ws.predictColumnTypes(dataRowIdx, len(header))
	fmt.Println(columnTypes)

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
func (ws *Worksheet) predictColumnTypes(dataStart int, columnCnt int) []string {
	row := ws.Sheet.Rows[dataStart]
	columnTypes := make([]string, len(row.Cells))

	// wee need the header row and the column count for this

	// i loops through each column once
	// j will loop through each value in that column
	for i := 0; i < columnCnt; i++ {
		// get info about this column
		// ? the idx == i
		// startingCell := ws.Sheet.Rows[dataStart].Cells[i]
		cp := &ColumnPrediction{}
		// cp := ws.determineColumnType(dataStart, i, "", startingCell, cp)
		fmt.Println(cp)
	}

	// ??: maybe we consider some kind of buffer for the columnTypes.... meaning if a max string length in a column is 200... should we set the max to 200+buffersize
	// ??: pretend buffer size is 55... in our example we would set the column type to varchar(255) so as to accomidate for the max size plus future sizes...
	// ??: same idea could apply to integers. Also nobody uses smallint. Also we sould have some detection to see if the value is a bool (always 0 or 1).
	// }
	return columnTypes
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
func (ws *Worksheet) determineColumnType(row, col int, prevType string, cell Cell, precision int, dups *map[string]bool, cp *ColumnPrediction) *ColumnPrediction {
	if cp.PrimaryKey || cp.NotNull {
		// get value
		value := ws.getValue(&cell)

		if (*dups)[value] == true {
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

	if prevType == "string" {
		// the type must remain a string

		// get precision

		// if precision is bigger than before, replace
		// else just use old precision

	}

	// ??? HOWWWWW?????

	if prevType == "boolean" {
		if cellType == prevType {
			// remains the same
		} else {
			// it's different
		}
	}

	if prevType == "float64" {
		if cellType == prevType {
			// remains the same
		} else {
			// it's different
		}
	}

	if prevType == "integer" {
		if cellType == prevType {
			// remains the same
		} else {
			// it's different
		}
	}

	// return a default case...??
}
