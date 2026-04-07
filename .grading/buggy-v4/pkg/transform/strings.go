package transform

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a snake_case or kebab-case string to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})

	if len(parts) == 0 {
		return ""
	}

	result := strings.ToLower(parts[0])
	for _, part := range parts[1:] {
		if len(part) == 0 {
			continue
		}
		result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return result
}

// ToSnakeCase converts a camelCase or PascalCase string to snake_case.
// BUG: produces double underscores for consecutive uppercase like "HTTPServer" -> "h_t_t_p__server"
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Slugify converts a string into a URL-friendly slug.
// FIXED: handles Unicode by transliterating common characters.
func Slugify(s string) string {
	if s == "" {
		return ""
	}

	// Basic Unicode transliteration
	replacements := map[rune]string{
		'é': "e", 'è': "e", 'ê': "e", 'ë': "e",
		'á': "a", 'à': "a", 'â': "a", 'ã': "a", 'ä': "a",
		'í': "i", 'ì': "i", 'î': "i", 'ï': "i",
		'ó': "o", 'ò': "o", 'ô': "o", 'õ': "o", 'ö': "o",
		'ú': "u", 'ù': "u", 'û': "u", 'ü': "u",
		'ñ': "n", 'ç': "c", 'ß': "ss",
	}

	s = strings.ToLower(s)

	var result strings.Builder
	prevDash := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			result.WriteRune(r)
			prevDash = false
		} else if rep, ok := replacements[r]; ok {
			result.WriteString(rep)
			prevDash = false
		} else if !prevDash && result.Len() > 0 {
			result.WriteRune('-')
			prevDash = true
		}
	}

	slug := result.String()
	slug = strings.TrimRight(slug, "-")
	return slug
}

// Abbreviate truncates a string to maxLen characters, appending "..." if truncated.
func Abbreviate(s string, maxLen int) string {
	if maxLen < 0 {
		maxLen = 0
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen] + "..."
	}
	return s[:maxLen-3] + "..."
}
