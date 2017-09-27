package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

type IOTData struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	BoardId     int           `json:"boardId"`
	Description string        `json:"description"`
	Temperature []float32     `json:"temperature"`
	DigitalIOA  uint8         `json:"digitalIOA"`
	DigitalIOB  uint8         `json:"digitalIOB"`
	Time        time.Time     `json:"timestamp"`
}

type Response struct {
	Status  string      `json:"status"`
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func IotcontrollerPostIOTData(w http.ResponseWriter, r *http.Request) {

	var response Response

	apikey := r.URL.Query().Get("apikey")
	if apikey != config.Apikey {
		response = Response{Status: "501", Result: "ko", Message: "ERROR you don't have correct permissions", Payload: nil}
	} else {
		r.ParseForm()

		server, err := config.GetServer("mongodb")
		if err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("Mongdb server %s", server.Host))

		// database setup and init
		session, err := mgo.Dial(server.Host)
		if err != nil {
			//panic(err)
			logger.Error(err.Error())
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)

		if err := session.DB("sampledb").Login(server.User, server.Password); err != nil {
			panic(err)
		}
		// Collection Rates
		c := session.DB("sampledb").C("temperature")

		logger.Info(fmt.Sprintf("Connected to mongodb host %s", server.Host))

		var iotdata IOTData
		body, err := ioutil.ReadAll(r.Body)
		logger.Info(fmt.Sprintf("Received json input %s", string(body[:])))
		errJson := json.Unmarshal(body, &iotdata)
		if errJson != nil {
			logger.Error(err.Error())
		}

		err = c.Insert(iotdata)
		if err != nil {
			panic(err)
		}
		response = Response{Status: "200", Result: "ok", Message: "INFO iotwebservice received payload ", Payload: nil}
	}

	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func ListIOTData(w http.ResponseWriter, r *http.Request) {

	var response Response
	apikey := r.URL.Query().Get("apikey")
	if apikey != config.Apikey {
		response = Response{Status: "501", Result: "ko", Message: "ERROR you don't have correct permissions", Payload: nil}
	} else {

		r.ParseForm()

		server, err := config.GetServer("mongodb")
		if err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("Mongdb server %s", server.Host))

		// database setup and init
		session, err := mgo.Dial(server.Host)
		if err != nil {
			//panic(err)
			logger.Error(err.Error())
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)

		if err := session.DB("sampledb").Login(server.User, server.Password); err != nil {
			panic(err)
		}

		var iotdata []IOTData
		c := session.DB("sampledb").C("temperature")
		logger.Info(fmt.Sprintf("Connected to mongodb host %s", server.Host))
		//query := c.Find(nil).Sort("timestamp").Limit(12)
		err = c.Find(nil).Sort("-timestamp").Limit(12).All(&iotdata)

		response = Response{Status: "200", Result: "ok", Message: "INFO iotwebservice received payload ", Payload: iotdata}
	}

	b, _ := json.Marshal(response)
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
