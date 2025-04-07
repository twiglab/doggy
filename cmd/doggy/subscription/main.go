package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/twiglab/doggy/holo"
	"resty.dev/v3"
)

type X struct {
	Enable  int
	PeerUrl string
}

func NewX(url string) *X {
	return &X{Enable: 1, PeerUrl: url}
}

var (
	addr string
	port int

	url string

	password string
	username string

	ssl int
)

func postUrl() string {
	if ssl == 0 {
		return fmt.Sprintf("http://%s:%d/SDCAPI/V1.0/Rest/RestClient", addr, port)
	}
	return fmt.Sprintf("https://%s:%d/SDCAPI/V1.0/Rest/RestClient", addr, port)
}

func init() {
	flag.IntVar(&port, "port", 80, "摄像头端口")
	flag.StringVar(&addr, "addr", "", "摄像头地址")

	flag.StringVar(&url, "url", "", "平台注册url")

	flag.StringVar(&username, "username", "", "用户名")
	flag.StringVar(&password, "password", "", "密码")

	flag.IntVar(&ssl, "ssl", 1, "是否使用ssl(0不启用)")
}

func main() {
	flag.Parse()

	var cr holo.CommonResponse

	c := resty.New().SetDigestAuth(username, password)
	defer c.Close()

	_, err := c.R().SetBody(NewX(url)).
		SetResult(&cr).
		SetError(&cr).
		Post(postUrl())

	if err != nil {
		log.Fatal(err)
	}

	log.Println(cr)
}
