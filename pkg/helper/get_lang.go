package helper

import (
	"strings"
)

func GetPrefixFromServerName(serverName string) string {
	splitStr := strings.Split(serverName, ".")
	if len(splitStr) > 0 {
		return splitStr[0]
	}

	return serverName
}
