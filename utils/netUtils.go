package utils

import (
	"strings"
	"fmt"
	"net/url"
	"strconv"
	"net/http"
	"time"
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

func HttpRequestWithLogger(handler http.HandlerFunc, name string) http.HandlerFunc {
	logger:= NewAppLogger("http-api", fmt.Sprintf("http[%s]", name))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start:= time.Now()
		handler.ServeHTTP(w, r)
		logger.Debugf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}