package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	Reporter    = "email"
	Environment = "-dev"
	//PORT        = ":29101"
)

func main() {
	godotenv.Load(".env")
	router := NewRouter()
	PORT := os.Getenv("PORT")
	log.Printf("Email Connector %s started on port: %s ", "T220106.9", PORT)

	log.Fatal(http.ListenAndServe(PORT, router))
}
