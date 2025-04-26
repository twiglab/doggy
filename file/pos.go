package file

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/twiglab/doggy/pf"
)

type CsvCameraPos struct {
	csvFileName string
	posMap      map[string]pf.CameraPos
}

func (r *CsvCameraPos) Get(sn string) (camera pf.CameraPos, ok bool) {
	camera, ok = r.posMap[sn]
	return
}

func (r *CsvCameraPos) Load() error {
	f, err := os.Open(r.csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	cr := csv.NewReader(f)
	cr.TrimLeadingSpace = true

	var l line

	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		l.reload(record)

		r.posMap[l.Pos(0)] = pf.CameraPos{
			SN:       l.Pos(0),
			Pos:      l.Pos(1),
			Floor:    l.Pos(2),
			Building: l.Pos(3),
			Area:     l.Pos(4),
		}
	}
	return nil
}
