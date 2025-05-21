package pf

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/twiglab/doggy/holo"
)

type cameraData struct {
	sn       string
	uuid     string
	deviceID string

	user string
	pwd  string
}

func buildCamera(rows []string) cameraData {
	return cameraData{
		sn:       rows[0],
		uuid:     rows[1],
		deviceID: rows[2],
	}
}

type CsvCameraDB struct {
	User string
	Pwd  string

	csvFile    string
	deviceConf map[string]cameraData

	strict bool
}

func NewCsvCameraDB(csvFile, user, pwd string) *CsvCameraDB {
	return &CsvCameraDB{
		User:       user,
		Pwd:        pwd,
		csvFile:    csvFile,
		deviceConf: make(map[string]cameraData),
	}
}

func (r *CsvCameraDB) Load(ctx context.Context) error {
	clear(r.deviceConf)

	if r.csvFile == "" {
		slog.InfoContext(ctx, "cameradb load", "csvfile", "<no csv file>")
		return nil
	}

	file, err := os.Open(r.csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	for {
		rows, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		device := buildCamera(rows)
		r.deviceConf[device.sn] = device
	}

	slog.InfoContext(ctx, "cameradb", slog.Int("size", len(r.deviceConf)))
	return nil
}

func (r *CsvCameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	if dev, ok := r.deviceConf[data.SerialNumber]; ok {
		device, err := holo.OpenDevice(data.IpAddr, r.User, r.Pwd)
		if err != nil {
			return nil, err
		}

		device.UUID = dev.uuid
		device.DeviceID = dev.deviceID
		device.SN = dev.sn

		return device, nil
	}

	if r.strict {
		return nil, fmt.Errorf("sn: %s not found", data.SerialNumber)
	}
	return holo.OpenDevice(data.IpAddr, r.User, r.Pwd)
}
