package lang

import (
	"strconv"
	"strings"
)

func FirstToUpper(str string) string {
	b := strings.ToUpper(str[:1]) + str[1:]
	return b
}

func ParseValueType(str string) interface{} {
	b, err := strconv.ParseBool(strings.ToLower(str))
	if err != nil {
		return str
	}
	return b
}
