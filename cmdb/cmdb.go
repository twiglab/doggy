package cmdb

import (
	"encoding/json/v2"
	"errors"
	"os"
)

type UserDB interface {
	ChannelData(cameraID, channelID string) (ChannelUserData, error)
}

type ChannelUserData struct {
	UUID string
	Code string

	X string
	Y string
	Z string
}

type CameraUserData struct {
	SN       string
	Channels []ChannelUserData
}

type UserData struct {
	CameraData []CameraUserData
}

func FormJson(file string) (u *UserData, err error) {
	var f *os.File
	u = &UserData{}

	f, err = os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.UnmarshalRead(f, u)
	return
}

func (u *UserData) ChannelData(cameraID, channelID string) (ch ChannelUserData, err error) {
	err = errors.New("not found")
	for _, c := range u.CameraData {
		if c.SN == cameraID {
			for _, ch = range c.Channels {
				if ch.UUID == channelID {
					err = nil
				}
			}
		}
	}

	return
}
