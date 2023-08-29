package xutil

import "strings"

func ContainsString(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func ContainsStringIgnoreCase(slice []string, target string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, target) {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, target int) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func UniqueMergeStringSlice(s1, s2 []string) []string {
	m := make(map[string]struct{}, len(s1)+len(s2))
	for _, s := range s1 {
		m[s] = struct{}{}
	}
	for _, s := range s2 {
		m[s] = struct{}{}
	}

	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}

	return r
}
