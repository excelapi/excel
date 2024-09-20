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
		fmt.Printf("Row %d has %d cells\n", i, len(row.Cells))

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
	for i := 0; i < 1; i++ {
		// get info about this column
		// ? the idx == i
		for j := dataStart; j < 5; j++ {
			// ? determine cell type
			cell := ws.Sheet.Rows[j].Cells[i]
			fmt.Println(j, i, cell.dataType(), ws.getValue(&cell))
			// ? determine maximum size of the datatype
			// ? determine uniqueness
			// ? determine empty values
		}
		// ? is nullable?
		// ? is primary key eligible
	}
	// for i, cell := range row.Cells {
	// if cell.Type == "s" {
	// todo: measure size of string
	// columnTypes[i] = "string"
	// } else {
	// todo: make sure that they don't have numbers higher than the maximum allowed for ints
	// ?: postgres doesn't allow for unsigned ints
	// ?: smallint 2 bytes = -32768 to +32767
	// ?: integer 4 bytes = -2147483648 to +2147483647
	// ?: bigint 8 bytes = -9223372036854775808 to +9223372036854775807
	// anything bigger should return an error
	// should default to the smallest possible option.
	// columnTypes[i] = "int"
	// }

	// ??: maybe we consider some kind of buffer for the columnTypes.... meaning if a max string length in a column is 200... should we set the max to 200+buffersize
	// ??: pretend buffer size is 55... in our example we would set the column type to varchar(255) so as to accomidate for the max size plus future sizes...
	// ??: same idea could apply to integers. Also nobody uses smallint. Also we sould have some detection to see if the value is a bool (always 0 or 1).
	// }
	return columnTypes
}

func (c *Cell) dataType() string {
	if c.Type == "s" {
		return "string"
	} else {
		if d, err := strconv.ParseFloat(c.Value, 64); err != nil {
			if math.Mod(d, 1) == 0 {
				return "int"
			}
			return "float64"
		}

		return "string"
	}
}
