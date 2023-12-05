package utils

import (
	"sort"
)

func ContainsString(sortedSlice []string, target string) bool {
	index := sort.SearchStrings(sortedSlice, target)
	return index < len(sortedSlice) && sortedSlice[index] == target
}
