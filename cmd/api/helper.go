package main

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/*
we are going to create three functions
1. To read JSON
2. To write JSON
3. To generate a JSON error response
*/

