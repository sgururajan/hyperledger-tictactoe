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

func HandleServerError(w http.ResponseWriter, additionalMsg string , err error) {
	w.WriteHeader(http.StatusInternalServerError)
	msg:= fmt.Sprintf("%s.\nErr: %v", additionalMsg, err)
	w.Write([]byte(msg))
}
