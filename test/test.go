package tests

import (
	"testing"
)

// integration testing.
// write function that uses this package as if it was a seperate application

func TestOpen(t *testing.T) {
	// file := file.Open("/Users/joshuablackhurst/Desktop/TestExcel.xlsx", "Sheet1")

	// file.ReadAll()
}

func TestPredictTable(t *testing.T) {
	ws := file.Open("test.xlsx", "Sheet1")
	prediction, err := ws.PredictTable()
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	if prediction.PrimaryKey == "" {
		t.Fatalf("primary key expected")
	}
}
