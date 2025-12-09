package be

const (
	TAOS  = "taos"
	MQTT  = "mqtt"
	MQTT5 = "mqtt5"
	NONE  = "none"
)

func HasHuman(in, out int) bool {
	return in != 0 || out != 0
}

func HasCount(count int) bool {
	return count != 0
}
