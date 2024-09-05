package main

import (
	"excel/internal/file"
	"fmt"
)

func main() {
	worksheet := file.Open("test.xlsx", "sheet1") // hard coded this will change later
	prediction, err := worksheet.PredictTable()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Table Prediction: %+v\n", prediction) // just a crappy log to see what we get
}
