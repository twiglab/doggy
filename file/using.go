package file

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/twiglab/doggy/pf"
)

type CsvCameraUsing struct {
	csvFileName string
	cameraMap   map[string]pf.CameraUsing
}

func (r *CsvCameraUsing) Load() error {
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

		u := pf.CameraUsing{
			SN:    l.Pos(0),
			UUID:  l.Pos(1),
			AlgID: l.Pos(2),
			Name:  l.Pos(3),
		}

		r.cameraMap[u.BK] = u
	}
	return nil
}

func NewBK(uuid, algid string) string {
	return algid + "@" + uuid
}
