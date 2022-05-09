package service

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	common "gitlab.com/dhf0820/ids_model/common"
	"net/http"
	"time"
	//logging "gitlab.com/dhf0820/ids_model/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckHealth() (status int, message string) {
	return 200, "OK"
}

type Message struct {
	DeliveryId primitive.ObjectID `json:"delivery_id"`
	LogType    string             `json:"log_type"`
	Status     string             `json:"status"`
	Reporter   string             `json:"reporter"`
	Message    string             `json:"message"`
	Timestamp  time.Time          `json:"timestamp"`
	Nanotime   int64              `json:"nanotime"`
}

func LogMessage(payload *common.Payload, logType string, status string, message string, url string) {
	var msg Message
	return
	msg.DeliveryId = payload.DelvPayload.ID
	msg.LogType = logType
	msg.Status = status
	msg.Reporter = Reporter
	msg.Message = message
	msg.Timestamp = time.Now().UTC()
	msg.Nanotime = time.Now().UnixNano()

	bstr, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling log message into json:", err.Error())
		return
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bstr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error sending log message:", err.Error())
		return
	}
	defer resp.Body.Close()
	return
}
