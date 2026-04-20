package utils

import "strings"

func SplitRoleKeys(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, item := range parts {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func JoinRoleKeys(keys []string) string {
	return strings.Join(SplitRoleKeys(strings.Join(keys, ",")), ",")
}

func HasRole(raw, target string) bool {
	for _, roleKey := range SplitRoleKeys(raw) {
		if roleKey == target {
			return true
		}
	}
	return false
}

func MergeDataScope(current, candidate string) string {
	rank := map[string]int{
		"self": 0,
		"dept": 1,
		"all":  2,
	}
	curRank, ok := rank[current]
	if !ok {
		curRank = 0
	}
	candRank, ok := rank[candidate]
	if !ok {
		candRank = 0
	}
	if candRank > curRank {
		return candidate
	}
	return current
}
