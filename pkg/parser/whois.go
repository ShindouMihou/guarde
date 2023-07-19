package parser

import (
	"guarde/pkg/utils"
	"strings"
)

func WhoIs(text string) map[string][]string {
	properties := make(map[string][]string)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if !strings.Contains(line, ":") {
			continue
		}
		arr := strings.SplitN(line, ":", 2)
		key, value := arr[0], arr[1]
		if utils.HasPrefixStr(key, "%") || utils.HasPrefixStr(key, "#") {
			continue
		}

		key = strings.TrimSpace(key)
		key = strings.ToLower(key)

		value = strings.TrimSpace(value)
		value = strings.ToLower(value)

		if _, ok := properties[key]; !ok {
			properties[key] = []string{value}
		} else {
			properties[key] = append(properties[key], value)
		}
	}
	return properties
}
