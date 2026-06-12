package pakasir

import "regexp"

var nonURLSafe = regexp.MustCompile(`[^\w\-_.~]`)

func sanitizeUrlSafe(s string) string {
	return nonURLSafe.ReplaceAllString(s, "")
}
