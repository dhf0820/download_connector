package service

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/dhf0820/ids_model/logging"

	// "gopkg.in/mgo.v2/bson"
	//"io/ioutil"
	"log"
	"net/http"

	common "gitlab.com/dhf0820/ids_model/common"
	// "strconv"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/gorilla/mux"
)

var LoggingUrl string
var StatusUrl string

func processPayload(res *http.Request) (*common.Payload, error) {
	var payload common.Payload

	//decoder := json.NewDecoder(res.Body).Decode(&payload)
	decoder := json.NewDecoder(res.Body).Decode(&payload)
	if decoder != nil {
		err := fmt.Errorf("Payload decode error: %s", decoder.Error())
		return nil, err
	}
	fmt.Printf("ProcessPayload:55 -- %s\n", spew.Sdump(payload))
	fmt.Printf("\n##CallBackMessageURL:30 -- %s\n",
		payload.DelvPayload.CallBackLogMessageUrl)
	msg := fmt.Sprintf("Delivery: %s received. Processing...")
	logging.LogMessage(&payload, "log", "Success", msg, "EmailConnector")

	return &payload, nil
}

func deliveryHandler(w http.ResponseWriter, r *http.Request) {
	resp := GenericResponse{}
	log.Println("DoDelivery request received")

	payload, err := processPayload(r)
	//reportingUrls := payload.Config.ReportingUrls
	LoggingUrl = payload.DelvPayload.CallBackLogMessageUrl //common.GetDataByName(reportingUrls, "logs")
	StatusUrl = payload.DelvPayload.CallBackStatusUrl      //common.GetDataByName(reportingUrls, "status" )

	fmt.Printf("LoggingURL: 49 -- [%s]\n", LoggingUrl)
	fmt.Printf("StatusURL: 50 -- [%s]\n", StatusUrl)
	// common.SendStatusReport("info","received", "delivery Payload",
	// 	payload.DelvPayload.CallBackStatusUrl)
	errCount := ValidatePayload(payload)
	if errCount > 0 {
		msg := fmt.Sprintf("A total of %d errors were found. Unprocessable!", errCount)
		logging.LogMessage(payload, "message", "error", msg, "EmailConnector")

		resp.Status = 400
		resp.Message = msg
		writeGenericResponse(w, resp.Status, resp.Message)
	} else {
		ProcessEmailPayload(payload)
		resp.Status = 200
		resp.Message = "Email accepted the request"
	}
	//TODO: This should be a go call ad return status of 200 Accepted
	//fmt.Printf("deliveryHandler Payload: %s\n", spew.Sdump(payload))
	// Process the email and send it
	//err = ProcessEmailPayload(payload)
	err = writeGenericResponse(w, resp.Status, resp.Message)
	if err != nil {
		log.Println("Error writing response: ", err.Error())
		return
	}
	return
}
