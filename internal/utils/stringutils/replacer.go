package stringutils

import "strings"

// ReplaceFunc заменяет все вхождения, соответствующие pattern (где %s - любая подстрока),
// используя функцию replacer для генерации замены для каждого совпадения
func ReplaceFunc(s, pattern string, replacer func(match string) string) string {
	parts := strings.Split(pattern, "%s")
	if len(parts) != 2 {
		return s
	}
	prefix, suffix := parts[0], parts[1]

	var result strings.Builder
	start := 0

	for {
		prefixPos := strings.Index(s[start:], prefix)
		if prefixPos == -1 {
			break
		}
		prefixPos += start

		afterPrefix := prefixPos + len(prefix)
		suffixPos := strings.Index(s[afterPrefix:], suffix)
		if suffixPos == -1 {
			break
		}
		suffixPos += afterPrefix

		match := s[afterPrefix:suffixPos]

		result.WriteString(s[start:prefixPos])

		replacement := replacer(match)
		result.WriteString(replacement)

		start = suffixPos + len(suffix)
	}

	result.WriteString(s[start:])

	return result.String()
}
