package file

import (
	"fmt"
	"strconv"
)

func (ws *Worksheet) PredictTable() (*TablePrediction, error) {
	headerIdx, err := ws.headerRowIndex()
	if err != nil {
		return nil, err
	}

	columnNames, err := ws.suggestHeader(headerIdx)
	if err != nil {
		return nil, err
	}

	dataRowIdx := headerIdx + 1
	columnTypes := ws.predictColumnTypes(dataRowIdx)

	predictedPK := columnNames[0]

	return &TablePrediction{
		ColumnNames: columnNames,
		ColumnTypes: columnTypes,
		PrimaryKey:  predictedPK,
	}, nil
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
func (ws *Worksheet) suggestHeader(idx int) ([]string, error) {
	header := ws.Sheet.Rows[idx]
	headerStr := []string{}

	for _, cell := range header.Cells {
		if cell.Type == "s" {
			idx, _ := strconv.Atoi(cell.Value)
			headerStr = append(headerStr, ws.SS.StringItems[idx].T)
		} else {
			return nil, fmt.Errorf("all header names must be strings")
		}
	}
	return headerStr, nil
}

// predict the types based
func (ws *Worksheet) predictColumnTypes(rowIdx int) []string {
	row := ws.Sheet.Rows[rowIdx]
	columnTypes := make([]string, len(row.Cells))

	for i, cell := range row.Cells {
		if cell.Type == "s" {
			// todo: measure size of string
			columnTypes[i] = "string"
		} else {
			// todo: make sure that they don't have numbers higher than the maximum allowed for ints
			// ?: postgres doesn't allow for unsigned ints
			// ?: smallint 2 bytes = -32768 to +32767
			// ?: integer 4 bytes = -2147483648 to +2147483647
			// ?: bigint 8 bytes = -9223372036854775808 to +9223372036854775807
			// anything bigger should return an error
			// should default to the smallest possible option.
			columnTypes[i] = "int"
		}

		// ??: maybe we consider some kind of buffer for the columnTypes.... meaning if a max string length in a column is 200... should we set the max to 200+buffersize
		// ??: pretend buffer size is 55... in our example we would set the column type to varchar(255) so as to accomidate for the max size plus future sizes...
		// ??: same idea could apply to integers. Also nobody uses smallint. Also we sould have some detection to see if the value is a bool (always 0 or 1).kj
	}
	return columnTypes
}
