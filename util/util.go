package util

import (
	"encoding/json"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Jsonp(obj interface{}, callback string) string {
	b, _ := json.Marshal(obj)
	return callback + "(" + string(b) + ")"
}
