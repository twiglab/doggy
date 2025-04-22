package file

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/twiglab/doggy/pf"
)

type CsvCameraUpload struct {
	csvFileName string
	cameraMap   map[string]pf.CameraUplaod
}

func (r *CsvCameraUpload) Get(sn string) (camera pf.CameraUplaod, ok bool) {
	camera, ok = r.cameraMap[sn]
	return
}

func (r *CsvCameraUpload) Load() error {
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

		r.cameraMap[l.Pos(0)] = pf.CameraUplaod{
			SN:     l.Pos(0),
			IpAddr: l.Pos(1),
			Last:   time.Now(),
		}
	}
	return nil
}
