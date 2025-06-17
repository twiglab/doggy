package d3252

import (
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/hx"
	"github.com/twiglab/doggy/job"
)

/*
00000000-0000-0000-0000-000000000000
ffffffff-ffff-ffff-ffff-ffffffffffff
*/
const (
	uuid     = "00000000-0000-0000-0000-000000000000"
	deviceID = "1234567890"
)

func rnd() int {
	return rand.IntN(5)
}

type Camera struct {
	DeviceAutoRegisterData holo.DeviceAutoRegisterData
	IDList                 holo.DeviceIDList
	SubMap                 map[string]holo.SubscriptionReq

	// ---------------------------------------------
	isAutoReg bool
}

var camera = &Camera{
	DeviceAutoRegisterData: holo.DeviceAutoRegisterData{
		SerialNumber: "ABCDEFGHIJKLMN",
		IpAddr:       "127.0.0.1:10007",
		DeviceName:   "kake SDC",
		Manufacturer: "fake",
		DeviceType:   "fake type",
		ChannelInfo:  []holo.Channel{{ChannelID: 101, UUID: uuid, DeviceID: deviceID}},
		DeviceVersion: holo.DeviceVersionData{
			Software: holo.SDC_11_0_0_SPC300,
		},
	},

	IDList: holo.DeviceIDList{IDs: []holo.DeviceID{
		{
			UUID:     uuid,
			DeviceID: deviceID,
		},
	}},
	SubMap: make(map[string]holo.SubscriptionReq),
}

var client = req.C().EnableInsecureSkipVerify()

func d3252() {
	cron, err := job.NewCron()
	if err != nil {
		log.Fatal(err)
	}

	cron.AddDurationFunc(5*time.Second, func() {
		if !camera.isAutoReg {
			var resp holo.CommonResponse

			_, err := client.R().
				SetBody(camera.DeviceAutoRegisterData).
				SetSuccessResult(&resp).
				SetErrorResult(&resp).
				Put("https://127.0.0.1:10005/pf/nat")

			if err != nil {
				log.Println("auto reg failet", err)
				return
			}
			if err = resp.Err(); err != nil {
				return
			}
			log.Println("auto reg ok")
			if !bBug {
				camera.isAutoReg = true
			}
		}
	})

	cron.AddDurationFunc(5*time.Second, func() {
		if !bEnableDensity {
			return
		}
		var resp holo.CommonResponse
		data := holo.MetadataObjectUpload{
			MetadataObject: holo.MetadataObject{
				Common: holo.Common{
					UUID:     camera.IDList.IDs[0].UUID,
					DeviceID: camera.IDList.IDs[0].DeviceID,
				},
				TargetList: []holo.HumanMix{
					{
						TargetType: holo.HUMMAN_DENSITY,
						HumanCount: rnd(),
						AreaRatio:  rnd(),
					},
				},
			},
		}
		for _, v := range camera.SubMap {
			_, err := client.R().
				SetBody(data).
				SetSuccessResult(&resp).
				SetErrorResult(&resp).
				Post(v.MetadataURL)

			if err != nil {
				log.Println("-----", err)
				return
			}
		}
	})

	cron.AddDurationFunc(30*time.Second, func() {
		var resp holo.CommonResponse
		data := holo.MetadataObjectUpload{
			MetadataObject: holo.MetadataObject{
				Common: holo.Common{
					UUID:     camera.IDList.IDs[0].UUID,
					DeviceID: camera.IDList.IDs[0].DeviceID,
				},
				TargetList: []holo.HumanMix{
					{
						TargetType:    holo.HUMMAN_COUNT,
						HumanCountIn:  rnd(),
						HumanCountOut: rnd(),
						StartTime:     time.Now().Add(-30 * time.Second).UnixMilli(),
						EndTime:       time.Now().UnixMilli(),
						TimeZone:      0, // for testing
					},
				},
			},
		}
		for _, v := range camera.SubMap {
			_, err := client.R().
				SetBody(data).
				SetSuccessResult(&resp).
				SetErrorResult(&resp).
				Post(v.MetadataURL)

			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	cron.Start()

	mux := chi.NewMux()
	mux.Use(middleware.Logger, middleware.Recoverer)

	sdcapi := chi.NewRouter()
	sdcapi.Post("/V1.0/System/Reboot", func(w http.ResponseWriter, r *http.Request) {
		clear(camera.SubMap)
		camera.isAutoReg = false

		camera.IDList = holo.DeviceIDList{IDs: []holo.DeviceID{
			{UUID: uuid, DeviceID: deviceID},
		}}

		hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	})

	sdcapi.Get("/V1.0/Rest/DeviceID", func(w http.ResponseWriter, r *http.Request) {
		hx.JsonTo(http.StatusOK, camera.IDList, w)
	})

	sdcapi.Put("/V1.0/Rest/DeviceID", func(w http.ResponseWriter, r *http.Request) {
		var idList holo.DeviceIDList
		if err := hx.BindAndClose(r, &idList); err != nil {
			hx.JsonTo(http.StatusInternalServerError,
				holo.NewCommonResponse(r.URL.Path, -1, err.Error()), w)
			return
		}

		if len(idList.IDs) > 0 {
			camera.IDList = idList
		}

		hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	})

	sdcapi.Get("/V2.0/Metadata/Subscription", func(w http.ResponseWriter, r *http.Request) {
		var data holo.Subscripions

		for _, v := range camera.SubMap {
			data.Subscriptions = append(data.Subscriptions, holo.SubscriptionItem{ID: 0, MetadataURL: v.MetadataURL})
		}
		hx.JsonTo(http.StatusOK, data, w)
	})

	sdcapi.Post("/V2.0/Metadata/Subscription", func(w http.ResponseWriter, r *http.Request) {
		var data holo.SubscriptionReq
		if err := hx.BindAndClose(r, &data); err != nil {
			hx.JsonTo(http.StatusInternalServerError,
				holo.NewCommonResponse(r.URL.Path, -1, err.Error()), w)
			return
		}
		camera.SubMap[data.MetadataURL] = data

		hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	})
	sdcapi.Delete("/V2.0/Metadata/Subscription", func(w http.ResponseWriter, r *http.Request) {
		clear(camera.SubMap)
		hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	})

	mux.Mount("/SDCAPI", sdcapi)

	if err := http.ListenAndServeTLS(":10007", "repo/server.crt", "repo/server.key", mux); err != nil {
		log.Fatal(err)
	}
}

var D3252Cmd = &cobra.Command{
	Use:   "D3252",
	Short: "D3252模拟器",
	Long:  `D3252模拟器`,
	Run: func(cmd *cobra.Command, args []string) {
		d3252()
	},

	Example: "dcp fake D3252",
}
var bBug bool
var bEnableDensity bool

func init() {
	D3252Cmd.Flags().BoolVar(&bBug, "bug", false, "bug模式")
	D3252Cmd.Flags().BoolVar(&bEnableDensity, "enable-density", false, "打开人流密度上报")
}
