package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type StatusCode int

const (
	InternalServerErrMessage = "Internal server error"
)

func WriteBadRequest(w http.ResponseWriter, r *http.Request, userMessage string) {
	log.WithFields(log.Fields{
		"method":   r.Method,
		"url":      r.URL,
		"httpCode": http.StatusBadRequest,
	}).Warn(userMessage)

	WriteError(w, userMessage, http.StatusBadRequest)
}

func WriteServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(log.Fields{
		"method":   r.Method,
		"url":      r.URL,
		"httpCode": http.StatusInternalServerError,
	}).Errorf("\n%+v\n", err)

	WriteError(w, InternalServerErrMessage, http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, userMessage string, httpErrorCode StatusCode) {
	http.Error(w, userMessage, int(httpErrorCode))
}

func WriteJSON(w http.ResponseWriter, r *http.Request, payload interface{}) {
	w.Header().Set("content-type", "application/json")

	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		WriteServerError(w, r, err)
		return
	}

	WriteResponse(w, r, jsonResponse)
}

func WriteResponse(w io.Writer, r *http.Request, resBytes []byte) {
	bytesWritten, err := w.Write(resBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"method":   r.Method,
			"url":      r.URL,
			"resBytes": fmt.Sprintf("[% x]", resBytes),
		}).Errorf("Err writing response bytes. \n%+v\n", err)

		return
	}

	log.WithFields(log.Fields{
		"method":       r.Method,
		"url":          r.URL,
		"bytesWritten": bytesWritten,
		"httpCode":     http.StatusOK,
	}).Info("HTTP response sent.")
}
