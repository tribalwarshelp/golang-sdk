package sdk

import "strings"

func minifyString(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
