package taos

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/taosdata/driver-go/v3/taosSql"
	_ "github.com/taosdata/driver-go/v3/taosWS"
)

func X() {
	var taosDSN = "root:taosdata@ws(localhost:6041)/"
	taos, err := sql.Open("taosWS", taosDSN)
	if err != nil {
		log.Fatalln("Failed to connect to " + taosDSN + "; ErrMessage: " + err.Error())
	}
	fmt.Println("Connected to " + taosDSN + " successfully.")
	defer taos.Close()
}
