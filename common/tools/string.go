package tools

import "sort"

// StrIN 高效搜索 href: https://cloud.tencent.com/developer/article/1729114
func StrIN(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
