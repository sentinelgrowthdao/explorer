package utils

import (
	"fmt"
	"strings"

	tmstrings "github.com/tendermint/tendermint/libs/strings"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseQuerySort(allowed []string, v string) (d bson.D, err error) {
	if v == "" {
		return d, nil
	}

	if !tmstrings.StringInSlice(v, allowed) {
		return nil, fmt.Errorf("sort value must be one of %#v", allowed)
	}

	keys := strings.Split(v, ",")
	for i := 0; i < len(keys); i++ {
		key, order := keys[i], 1
		if key[0] == '-' {
			key, order = key[1:], -1
		}

		d = append(d, bson.E{Key: key, Value: order})
	}

	return d, nil
}
