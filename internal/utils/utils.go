package utils

import (
	"strconv"
	"strings"
)

func DeepGet(m map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var cur any = m

	for _, p := range parts {
		switch v := cur.(type) {
		case map[string]any:
			cur = v[p]
		case []any:
			// optional: support array access like "items.0.name"
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil, false
			}
			cur = v[idx]
		default:
			return nil, false
		}
	}

	return cur, true
}
