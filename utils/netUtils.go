package utils

import (
	"strings"
	"fmt"
	"net/url"
	"strconv"
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