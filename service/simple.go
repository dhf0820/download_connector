package service

import (
	"fmt"

	//"github.com/davecgh/go-spew/spew"
	"gitlab.com/dhf0820/ids_model/common"
	docMod "gitlab.com/dhf0820/ids_model/document"
	"gitlab.com/dhf0820/ids_model/logging"

	//"io/ioutil"
	//"log"
	//"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func Send(payload *common.Payload, dDocs []*docMod.DeliveryDocument) error {

	smtpServer := payload.Config.ConnectAddress.Address
	auth := payload.Config.ConnectAddress.Authorization
	userName := ""
	password := ""
	if auth != "" {
		auths := strings.Split(auth, ":")
		//fmt.Printf("Auth: %s\n", auths)
		userName = auths[0]
		password = auths[1]
	}
	fmt.Printf("Access:%s:%s\n", userName, password)
	fromFld, err := common.GetFieldByName(payload.Config.Fields, "from")
	if err != nil {
		fmt.Printf("From Error: %s\n", err.Error())
	}
	toFld, err := common.GetFieldByName(payload.Device.Fields, "to")
	if err != nil {
		fmt.Printf("toFld Error: %s\n", err.Error())
	}
	fmt.Printf("Sending from: %s  to: %s\n", fromFld.Default, toFld.Value)
	gm := gomail.NewMessage()
	gm.SetHeader("From", fromFld.Default)
	gm.SetHeader("To", toFld.Value)
	fmt.Printf("---Sending to: %s---\n", toFld.Value)
	gm.SetHeader("Subject", "#secure# Requested Secure Medical Records")

	meta := payload.DelvPayload.Meta
	dob := common.GetDataByName(meta, "dob")
	patientName := common.GetDataByName(meta, "patient_name")

	body := fmt.Sprintf("Requested Medical Records for Patient: <b>%s</b> -dob: <i>%s</i>!", patientName, dob)
	gm.SetBody("text/html", body)
	for _, dDoc := range dDocs {

		fmt.Printf("     Attaching file: [%s]  name: [%s]\n", dDoc.Description, dDoc.FileName)
		//fmt.Printf("attached Delivery: %s\n", spew.Sdump(dDoc))
		//fi, _ := os.Stat(dDoc.FileName)
		//fmt.Printf("Size of %s : %d\n\n", dDoc.FileName, fi.Size())
		gm.Attach(dDoc.FileName)
	}
	fmt.Printf("SMTP SERVER: %s\n", smtpServer)
	var port int
	hostInfo := strings.Split(smtpServer, ":")
	if len(hostInfo) > 1 {
		port, err = strconv.Atoi(hostInfo[1])
		if err != nil {
			port = 587
		}
	}
	smtpServer = hostInfo[0]
	d := gomail.NewDialer(smtpServer, port, userName, password)
	startTime := time.Now()
	if err := d.DialAndSend(gm); err != nil {
		_ = removeTempFile(dDocs)
		fmt.Printf("DialAndSend error: %v\n", err)
		return err
	}
	msg := fmt.Sprintf("Send email of %d documents took %f seconds", len(dDocs), time.Since(startTime).Seconds)
	logging.LogMessage(payload, "logs", "success", msg, "email_connector")
	return nil
}

func removeTempFile(dDocs []*docMod.DeliveryDocument) error {
	for _, dDoc := range dDocs {
		fmt.Printf("Remove %s\n", dDoc.FileName)
		//os.Remove(dDoc.FileName)
		//if err != nil {
		//	fmt.Printf("Error removing tempFile %s: %v\n", doc.FileName, err)
		//}
	}
	return nil
}
