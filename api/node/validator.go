package node

import (
	"github.com/sentinel-official/explorer/utils"
)

var (
	validators = map[string]func(req *RequestGetNodeStatistics) error{
		"": validateHistorical,
	}
)

func validateHistorical(req *RequestGetNodeStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}
