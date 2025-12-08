package mqttc

func pushTopic(uuid, typ string) string {
	return "dcp/" + uuid + "/" + typ
}
