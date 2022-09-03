package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

//Creating a struct type that handles all the actions from various microservices
type RequestPayload struct {
	Action         string      `json:"action"`
	Authentication AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	//write the payload data
	_ = app.writeJSON(w, http.StatusOK, payLoad)

}

//handle submssion handles submitting authenticated microservices
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	//read in a var of type request payload
	var requestPayload RequestPayload

	//read json body
	err := app.readJSON(w, r, requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("unable to read data"), http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Authentication)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

//authenticate auth microservice
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	//create some JSON and send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the auth service, create an instance of a client
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	//create a variable to read the response body into
	var jsonFromService jsonResponse

	//decode the json from the auth sevice
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}
