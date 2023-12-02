package utils

import (
	"sort"
)

func ContainsString(sortedSlice []string, target string) bool {
	index := sort.SearchStrings(sortedSlice, target)
	return index < len(sortedSlice) && sortedSlice[index] == target
}

func FindStringIndex(sortedSlice []string, target string) int {
	index := sort.SearchStrings(sortedSlice, target)
	if index < len(sortedSlice) && sortedSlice[index] == target {
		return index
	}
	return -1
}
