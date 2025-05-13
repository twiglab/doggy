package pf

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/twiglab/doggy/holo"
)

type deviceData struct {
	sn       string
	uuid     string
	deviceID string

	user string
	pwd  string
}

func buildDev(rows []string) deviceData {
	return deviceData{
		sn:       rows[0],
		uuid:     rows[1],
		deviceID: rows[2],
	}
}

type CsvResolve struct {
	User string
	Pwd  string

	csvFile    string
	deviceConf map[string]deviceData

	strict bool
}

func NewCsvConfDeviceResolve(csvFile, user, pwd string) *CsvResolve {
	return &CsvResolve{
		User:       user,
		Pwd:        pwd,
		csvFile:    csvFile,
		deviceConf: make(map[string]deviceData),
	}
}

func (r *CsvResolve) Init() error {
	clear(r.deviceConf)

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
		device := buildDev(rows)
		r.deviceConf[device.sn] = device
	}
	return nil
}

func (r *CsvResolve) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
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
