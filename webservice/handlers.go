package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func ProcessIOTData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
	logger.Info(fmt.Sprintf("Received json input %s", string(body[:])))
	errJson := json.Unmarshal(body, &f)
	if errJson != nil {
		logger.Error(err.Error())
	}

	//m := f.(map[string]interface{})
	var response Response
	response = Response{Status: "200", Result: "ok", Message: "INFO iotwebservice received payload ", Payload: nil}
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
