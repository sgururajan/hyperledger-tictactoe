package common

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func GetPortFromUrl(urlStr string) uint32 {
	var port uint32
	port = 0

	slashIndex := strings.Index(urlStr, "//")
	if slashIndex < 0 {
		urlStr = fmt.Sprintf("http://%s", urlStr)
	}

	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return port
	}

	portUInt64, err := strconv.ParseUint(urlObj.Port(), 10, 32)
	if err != nil {
		return port
	}

	return uint32(portUInt64)
}

func HasStringKey(dict map[string]interface{}, consName string) bool {
	for k, _ := range dict {
		if k == consName {
			return true
		}
	}

	return false
}

func HasConsortium(dict map[string][]string, consName string) bool {
	for k, _ := range dict {
		if k == consName {
			return true
		}
	}

	return false
}
