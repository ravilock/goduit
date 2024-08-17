package integrationtests

import "strings"

func MakeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	return strings.ReplaceAll(loweredTitle, " ", "-")
}
