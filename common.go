package fcgiserver

import (
	"net/http"

	"github.com/s2ks/fcgiserver/logger"
)

func LogRequest(r *http.Request) {
	logger.Infof("HTTP Method: %s", r.Method)
	logger.Infof("Protocol version: %s", r.Proto)
	logger.Infof("Header: %+v", r.Header)
	logger.Infof("Client: %s", r.RemoteAddr)
	logger.Infof("Content length: %d", r.ContentLength)
}

func LogBody(r *http.Request, l uint64) {
	body := make([]byte, l)

	n, err := r.Body.Read(body)

	if err != nil || n == 0 {
		return
	}

	logger.Infof("Body (%v byte(s)): %s", l, string(body))
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	logger.Errorf("Error while serving %s - %v", r.URL.Path, err)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	logger.Errorf("%s not found (404)", r.URL.Path)
}
