package backend

const (
	TAOS = "taos"
	MQTT = "mqtt"
	NONE = "none"
)

func HasHuman(in, out int) bool {
	return in != 0 || out != 0
}

func HasCount(count int) bool {
	return count != 0
}
