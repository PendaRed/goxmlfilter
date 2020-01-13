package filterservice

import (
	"strings"
	"time"
)

func debugTStamp() string {
	return time.Now().Format("15:04:05.000")
}

// path has form el1/el2/el3
// Rules is keyed off the tail of the path, eg el2/el3
// This sees if a rule matches the end of the path, and if so returns the array of values
// which matches the rule
func findMatchingRule(path string, rules map[string][]string) ([]string, bool) {
	for k, v := range rules {
		if strings.HasSuffix(path, k) {
			return v, true
		}
	}
	return nil, false
}
