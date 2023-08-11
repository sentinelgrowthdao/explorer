package utils

import (
	"encoding/json"
)

func MustMarshal(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(buf)
}

func MustMarshalIndent(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(buf)
}
