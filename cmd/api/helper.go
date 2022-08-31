package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/*
In this helper package we are going to create three functions
1. To read JSON
2. To write JSON
3. To generate a JSON error response
*/

//1. To read JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	//set a max value for read JSON data
	maxByt := 1048576 //one megabytes

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByt))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	//check if body has a single JSON value
	err = dec.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must have one single JSON value")
	}

	return nil
}

//2. To write JSON
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	//Marshal JSON
	byt, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//check if headers was included
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(byt)
	if err != nil {
		return err
	}

	return nil
}


//3. To generate a JSON error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	//default status code
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error  = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)

}