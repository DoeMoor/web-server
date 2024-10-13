package api

import "net/http"

type apiHandler struct {
}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request){}

