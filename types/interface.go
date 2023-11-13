package types

import (
	"encoding/json"
	"time"
)

func Int64FromInterface(v interface{}) int64 {
	if v == nil {
		return 0
	}

	return v.(int64)
}

func StringFromInterface(v interface{}) string {
	if v == nil {
		return ""
	}

	return v.(string)
}

func TimeFromInterface(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}

	return v.(time.Time)
}

func Uint64FromInterface(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	if v, ok := v.(int64); ok {
		return uint64(v)
	}

	return v.(uint64)
}

func BandwidthFromInterface(v interface{}) *Bandwidth {
	if v == nil {
		return NewBandwidth(nil)
	}

	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	var res Bandwidth
	if err := json.Unmarshal(buf, &res); err != nil {
		panic(err)
	}

	return &res
}
