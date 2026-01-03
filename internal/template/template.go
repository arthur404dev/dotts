package template

import (
	"regexp"
)

var placeholderRegex = regexp.MustCompile(`<<dotts:([a-z][a-z0-9_.]*)>>`)

func Apply(content string, values map[string]string) string {
	return placeholderRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatch := placeholderRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		key := submatch[1]
		if val, ok := values[key]; ok && val != "" {
			return val
		}

		return match
	})
}

func ApplyBytes(content []byte, values map[string]string) []byte {
	return []byte(Apply(string(content), values))
}

func HasPlaceholders(content string) bool {
	return placeholderRegex.MatchString(content)
}

func HasPlaceholdersBytes(content []byte) bool {
	return placeholderRegex.Match(content)
}

func ExtractPlaceholders(content string) []string {
	matches := placeholderRegex.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var keys []string

	for _, match := range matches {
		if len(match) >= 2 {
			key := match[1]
			if !seen[key] {
				seen[key] = true
				keys = append(keys, key)
			}
		}
	}

	return keys
}

func GetMissingKeys(content string, values map[string]string) []string {
	keys := ExtractPlaceholders(content)

	var missing []string
	for _, key := range keys {
		if val, ok := values[key]; !ok || val == "" {
			missing = append(missing, key)
		}
	}

	return missing
}
