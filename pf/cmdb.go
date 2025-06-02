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

type cameraData struct {
	sn  string
	mac string

	uuid     string
	deviceID string

	user string
	pwd  string
}

func buildCamera(rows []string) cameraData {
	return cameraData{
		sn:       rows[0],
		mac:      rows[1],
		uuid:     rows[2],
		deviceID: rows[3],
	}
}

type CsvCameraDB struct {
	User string
	Pwd  string

	csvFile string
	snMap   map[string]cameraData
	uuidMap map[string]cameraData

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

func (r *CsvCameraDB) load(ctx context.Context) (map[string]cameraData, map[string]cameraData, error) {
	snMap := make(map[string]cameraData)
	uuidMap := make(map[string]cameraData)

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
		snMap[device.sn] = device
		uuidMap[device.uuid] = device
	}

	slog.InfoContext(ctx, "cameradb", slog.Int("size", len(r.snMap)))
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

func (r *CsvCameraDB) GetBySn(sn string) (data cameraData, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok = r.snMap[sn]
	return
}

func (r *CsvCameraDB) GetByUUID(uuid string) (data cameraData, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok = r.uuidMap[uuid]
	return
}

func (r *CsvCameraDB) SetTTL(sn string, time int64) {
	r.liveMap[sn] = time
}

func (r *CsvCameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.SerialNumber, data.IpAddr, r.User, r.Pwd)
}
