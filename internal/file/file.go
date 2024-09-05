package file

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
)

type sharedStrings struct {
	StringItems []StringItem `xml:"si"`
}

type StringItem struct {
	T string `xml:"t"`
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

type TablePrediction struct {
	ColumnNames []string
	ColumnTypes []string // ?: this array should be:  [int, varchar(13), varchar(255), datetime]
	PrimaryKey  string
}

func parseSharedStrings(f *zip.File, ws *Worksheet) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	return decoder.Decode(&ws.SS)
}

func Open(filepath, sheetName string) *Worksheet {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer r.Close()

	var ws Worksheet

	fullSheetName := fmt.Sprintf("xl/worksheets/%s.xml", sheetName)
	//fullSheetName := "xl/worksheets/sheet1.xml"

	for _, f := range r.File {
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

	return &ws
}
