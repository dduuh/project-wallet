package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var ErrHTTPMethod = errors.New("incorrect HTTP method")

type Map map[string]interface{}

func response(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return
	}

	resp, err := json.Marshal(message)
	if err != nil {
		logrus.Panicf("JSON Marshal error: %v\n", err)
	}

	if _, err = w.Write(resp); err != nil {
		logrus.Panicf("Write HTTP error: %v\n", err)
	}
}
