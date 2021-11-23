package util

import "strings"

func RemoveRepetitionIntSlice(targetList []int) (result []int) {
	if len(targetList) == 0 {
		return
	}
	repetitionMap := make(map[int]struct{}, 0)
	for _, item := range targetList {
		_, ok := repetitionMap[item]
		if !ok {
			result = append(result, item)
			repetitionMap[item] = struct{}{}
		}
	}
	return
}

func RemoveRepetitionStringSlice(targetList []string) (result []string) {
	if len(targetList) == 0 {
		return
	}
	repetitionMap := make(map[string]struct{}, 0)
	for _, item := range targetList {
		_, ok := repetitionMap[item]
		if !ok {
			result = append(result, item)
			repetitionMap[item] = struct{}{}
		}
	}
	return
}

//RemoveSliceEmpty *******删除slice里面的空格和去除字符的前后的空格*******//
func RemoveSliceEmpty(arg []string) []string {
	var keys []string
	for _, key := range arg {
		if key != "" || strings.Trim(key, " ") != "" {
			keys = append(keys, key)
		}
	}
	return keys
}
