package helpers

import "strings"

func ToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if strings.Contains(strings.ToLower(words[i]), "id") {
			words[i] = "ID"
			break
		}
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, "")
}