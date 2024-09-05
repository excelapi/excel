package file

import (
	"archive/zip"
	"encoding/xml"
)

func parseSheet(f *zip.File, ws *Worksheet) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	return decoder.Decode(&ws)
}
