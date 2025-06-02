package pf

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/twiglab/doggy/holo"
)

type CameraData struct {
	SN  string
	Mac string

	UUID     string
	DeviceID string

	User string
	Pwd  string
}

func buildCamera(rows []string) CameraData {
	return CameraData{
		SN:       rows[0],
		Mac:      rows[1],
		UUID:     rows[2],
		DeviceID: rows[3],
	}
}

type CsvCameraDB struct {
	User string
	Pwd  string

	csvFile string
	snMap   map[string]CameraData
	uuidMap map[string]CameraData

	liveMap map[string]int64

	mu sync.RWMutex
}

func NewCsvCameraDB(csvFile, user, pwd string) *CsvCameraDB {
	return &CsvCameraDB{
		User:    user,
		Pwd:     pwd,
		csvFile: csvFile,
		liveMap: make(map[string]int64),
	}
}

func (r *CsvCameraDB) load(ctx context.Context) (map[string]CameraData, map[string]CameraData, error) {
	snMap := make(map[string]CameraData)
	uuidMap := make(map[string]CameraData)

	if r.csvFile == "" {
		slog.InfoContext(ctx, "no csv file")
		return snMap, uuidMap, nil
	}

	file, err := os.Open(r.csvFile)
	if err != nil {
		return snMap, uuidMap, err
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
			return snMap, uuidMap, err
		}
		device := buildCamera(rows)
		snMap[device.SN] = device
		uuidMap[device.UUID] = device
	}

	slog.InfoContext(ctx, "cameradb", slog.Int("size", len(snMap)))
	return snMap, uuidMap, nil
}

func (r *CsvCameraDB) Load(ctx context.Context) error {

	snMap, uuidMap, err := r.load(ctx)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.snMap = snMap
	r.uuidMap = uuidMap

	return nil
}

func (r *CsvCameraDB) GetBySn(sn string) (data CameraData) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data = r.snMap[sn]
	return
}

func (r *CsvCameraDB) GetByUUID(uuid string) (data CameraData) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data = r.uuidMap[uuid]
	return
}

func (r *CsvCameraDB) SetTTL(uuid string, time int64) {
	r.liveMap[uuid] = time
}

func (r *CsvCameraDB) GetTTL(uuid string) int64 {
	i, b := r.liveMap[uuid]
	if b {
		return i
	}
	return 0
}

func (r *CsvCameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.SerialNumber, data.IpAddr, r.User, r.Pwd)
}
