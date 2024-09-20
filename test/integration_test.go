package test

import (
	"testing"

	"github.com/excelapi/excel/pkg/xlsx"
)

// integration testing.
// write function that uses this package as if it was a seperate application

func TestOpen(t *testing.T) {
	file := xlsx.Open("/Users/joshuablackhurst/source/repos/personal/excel/raw/TestExcel.xlsx", "sheet1")

	file.ReadAll()
}

// func TestWriteSql(t *testing.T) {
// 	file := xlsx.Open("/Users/joshuablackhurst/Desktop/TestExcel.xlsx", "sheet1")

// 	// pull in table-prediction.json
// 	jsn, _ := os.Open("./json/table-prediction.json")

// 	// get the bytes
// 	bytes, _ := io.ReadAll(jsn)

// 	// create table prediction
// 	var tp xlsx.TablePrediction
// 	json.Unmarshal(bytes, &tp)

// 	// send tp to WriteSQL()
// 	file.WriteSQL(&tp)
// }
