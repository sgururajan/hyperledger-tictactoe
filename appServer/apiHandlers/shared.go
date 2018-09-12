package apiHandlers

import (
	"net/http"
	"fmt"
)

type Route struct {
	Method      string
	Name        string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func HandleServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	msg:= fmt.Sprintf("msg: %s\nexception: %#v", err)
	w.Write([]byte(msg))
}
