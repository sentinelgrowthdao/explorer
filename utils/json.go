package utils

import (
	"encoding/json"
)

func MustMarshal(v interface{}) []byte {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return buf
}

func MustMarshalToString(v interface{}) string {
	return string(MustMarshal(v))
}

func MustMarshalIndent(v interface{}) []byte {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return buf
}

func MustMarshalIndentToString(v interface{}) string {
	return string(MustMarshalIndent(v))
}
