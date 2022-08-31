package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	//write the payload data
	_ = app.writeJSON(w, http.StatusOK, payLoad)

}
