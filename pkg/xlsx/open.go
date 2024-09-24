package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"strconv"
)

type sharedStrings struct {
	StringItems []StringItem `xml:"si"`
}

type StringItem struct {
	Text string `xml:"t"`
}

type Sheet struct {
	Rows []Row `xml:"row"`
}

type Row struct {
	Cells []Cell `xml:"c"`
}
type Cell struct {
	Type  string `xml:"t,attr"`
	Value string `xml:"v"`
}

type Worksheet struct {
	Sheet Sheet `xml:"sheetData"`
	SS    *sharedStrings
}

func Open(filepath, sheetName string) *Worksheet {

	// todo: validate filepath & sheetName

	r, err := zip.OpenReader(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer r.Close()

	var ws Worksheet

	// ? xlsx files are basically just a zip file - LAME
	for _, f := range r.File {
		fullSheetName := fmt.Sprintf("xl/worksheets/%v.xml", sheetName)

		if f.Name == fullSheetName {
			rc, err := f.Open()
			if err != nil {
				panic("can't open sheet: " + err.Error())
			}
			defer rc.Close()

			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&ws)
			if err != nil {
				panic("could not decode sheet: " + err.Error())
			}
		} else if f.Name == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err != nil {
				panic("can't open sharedStrings sheet: " + err.Error())
			}
			defer rc.Close()

			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&ws.SS)
			if err != nil {
				panic("could not decode sharedStrings sheet: " + err.Error())
			}
		}
	}

	// todo: filter out empty rows
	ws.filterEmptyRows()

	return &ws
}

<<<<<<< HEAD
func (ws *Worksheet) getValue(cell *Cell) string {
	var value string
	if cell.Type == "s" {
		idx, _ := strconv.Atoi(cell.Value)
		value = ws.SS.StringItems[idx].Text
	} else {
		// int, or something
		value = cell.Value
	}
	return value
}

// func (ws *Worksheet) headerRowIndex() (int, error) {
// 	data := ws.Sheet.Rows
// 	for i := 0; i < len(data); i++ {
// 		if len(data[i].Cells) > 0 {
// 			// potential header found
// 			return i, nil
// 		}
// 	}
// 	return 0, fmt.Errorf("unable to locate potential header row")
// }
=======
func (ws *Worksheet) isRowEmpty(row Row) bool {
	for _, cell := range row.Cells {
		if ws.getValue(&cell) != "" {
			return false
		}
	}
	return true
}
>>>>>>> origin/joshua/filter-empty-rows

func (ws *Worksheet) filterEmptyRows() {
	var nonEmptyRows []Row
	for _, row := range ws.Sheet.Rows {
		if !ws.isRowEmpty(row) {
			nonEmptyRows = append(nonEmptyRows, row)
		}
	}
	ws.Sheet.Rows = nonEmptyRows
}

<<<<<<< HEAD
// 	for _, cell := range header.Cells {
// 		if cell.Type == "s" {
// 			idx, _ := strconv.Atoi(cell.Value)
// 			headerStr = append(headerStr, ws.SS.StringItems[idx].T)
// 		} else {
// 			return []string{}, fmt.Errorf("all header names must be strings")
// 		}
// 	}

// 	return headerStr, nil
// }

// func (ws *Worksheet) getString(c *Cell) string {
// 	idx, _ := strconv.Atoi(c.Value)
// 	return ws.SS.StringItems[idx].Text
// }

// func (ws *Worksheet) ReadAll() {
// 	for _, row := range ws.Sheet.Rows {
// 		for _, cell := range row.Cells {
// 			var value string
// 			if cell.Type == "s" {
// 				value = ws.getString(&cell)
// 			} else {
// 				// int, or something
// 				value = cell.Value
// 			}
// 			fmt.Println(value)
// 		}
// 	}
// }
=======
func (ws *Worksheet) getValue(cell *Cell) string {
	var value string
	if cell.Type == "s" {
		idx, _ := strconv.Atoi(cell.Value)
		value = ws.SS.StringItems[idx].Text
	} else {
		value = cell.Value
	}
	return value
}

func (ws *Worksheet) ReadAll() {
	for _, row := range ws.Sheet.Rows {
		for _, cell := range row.Cells {
			var value string
			if cell.Type == "s" {
				value = ws.getValue(&cell)
			} else {
				// int, or something
				value = cell.Value
			}
			fmt.Println(value)
		}
	}

	fmt.Println("Rows: ", len(ws.Sheet.Rows))
	fmt.Println("First col length: ", len(ws.Sheet.Rows[0].Cells))
}
>>>>>>> origin/joshua/filter-empty-rows
