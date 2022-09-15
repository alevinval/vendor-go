package utils

import "sort"

func Union(a, b []string) []string {
	union := map[string]struct{}{}
	for i := range a {
		union[a[i]] = struct{}{}
	}
	for i := range b {
		union[b[i]] = struct{}{}
	}

	list := make([]string, 0, len(union))
	for key := range union {
		list = append(list, key)
	}
	sort.Strings(list)
	return list
}
