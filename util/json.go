package util

import "encoding/json"

func JsonStringToObject[T any](jsonStr string) T {
	var res T
	json.Unmarshal([]byte(jsonStr), &res)
	return res
}

func ObjectToJsonString(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}
