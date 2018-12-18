package utils

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func RespondIfError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func FailIfError(err error) {
	if err == nil {
		return
	}
	log.Panic(err.Error())
}
