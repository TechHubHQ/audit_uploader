package helpers

import "strings"

func RemoveSpaces(str string) string {
	result := strings.Join(strings.Fields(str), "")
	return result
}
