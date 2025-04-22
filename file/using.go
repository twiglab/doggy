package file

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/twiglab/doggy/mdm"
)

type CsvCameraUsing struct {
	csvFileName string
	cameraMap   map[string]mdm.CameraUsing
}

func (r *CsvCameraUsing) GetByBK(bk string) (u mdm.CameraUsing, ok bool) {
	u, ok = r.cameraMap[bk]
	return
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

		u := mdm.CameraUsing{
			SN:    l.Pos(0),
			UUID:  l.Pos(1),
			AlgID: l.Pos(2),
			Name:  l.Pos(3),
			Memo:  l.Pos(4),
			BK:    NewBK(l.Pos(2), l.Pos(3)),
		}

		r.cameraMap[u.BK] = u
	}
	return nil
}

func NewBK(uuid, algid string) string {
	return algid + "@" + uuid
}
