package service

import (
	"encoding/json"
	"log"
	"net/http"
	// "time"
)

type GenericResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func writeGenericResponse(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	switch status {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 403:
		w.WriteHeader(http.StatusForbidden)
	}
	resp := GenericResponse{Status: status, Message: message}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error marshaling json:", err.Error())
		return err
	}
	return nil
}

//standard log process to log messages to the core
// type Message struct {
// 	DeliveryId      bson.ObjectId `json:"delivery_id"`
// 	LogType         string        `json:"log_type"`
// 	Status          string        `json:"status"`
// 	Connector       string        `json:"connector"`
// 	Message         string        `json:"message"`
// 	Timestamp       time.Time     `json:"timestamp"`
// 	Nanotime        int64         `json:"nanotime"`
// }

func healthcheck(w http.ResponseWriter, r *http.Request) {
	status, message := CheckHealth()
	writeGenericResponse(w, status, message)
}
