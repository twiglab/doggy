package file

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/twiglab/doggy/pf"
)

type CsvCameraSetup struct {
	csvFileName string
	posMap      map[string]pf.CameraSetup
}

func (r *CsvCameraSetup) Get(sn string) (camera pf.CameraSetup, ok bool) {
	camera, ok = r.posMap[sn]
	return
}

func (r *CsvCameraSetup) Load() error {
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

		r.posMap[l.Pos(0)] = pf.CameraSetup{
			SN:       l.Pos(0),
			Pos:      l.Pos(1),
			Floor:    l.Pos(2),
			Building: l.Pos(3),
			Area:     l.Pos(4),
		}
	}
	return nil
}
