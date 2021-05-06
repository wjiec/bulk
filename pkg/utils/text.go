package utils

import "unicode"

// CamelCaseToHyphen converts from camel case form to hyphen separated form
// Example: HelloWorld -> hello-world
func CamelCaseToHyphen(s string) string {
	var out []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 && s[i-1] != '-' {
			out = append(out, '-')
		}
		out = append(out, unicode.ToLower(r))
	}

	return string(out)
}
