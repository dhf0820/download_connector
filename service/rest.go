package service

import (
	"bytes"
	common "gitlab.com/dhf0820/ids_model/common"
	docMod "gitlab.com/dhf0820/ids_model/document"
	"gitlab.com/dhf0820/ids_model/logging"
	"os"

	//"encoding/json"
	//"errors"
	"fmt"

	//"github.com/davecgh/go-spew/spew"

	// "io/ioutil"
	//"log"
	//logging "gitlab.com/dhf0820/ids_model/logging"
	// "github.com/gorilla/mux"
	"net/http"
	//"github.com/oleiade/reflections"
	"strconv"
	"strings"
)

var errCount = 0

func ValidatePayload(payload *common.Payload) int {
	status := &common.Status{
		State:   "Started",
		Comment: "Validating Payload",
	}

	common.SendStatus("Info", status, payload.DelvPayload.CallBackStatusUrl)
	errCount = 0
	ValidateConfig(payload)

	//Check for at least one document
	dp := payload.DelvPayload
	if len(dp.Documents) == 0 {
		logging.LogMessage(payload, "status", "error", "No Documents to deliver", "EmailConnector")
		errCount++
	}
	ValidateDevice(payload)
	//if err != nil {
	//	logging.LogMessage(payload, "basic", "error", err.Error(), payload.Config.LogURL, "EmailConnector")
	//	errCount++
	//}
	//if errCount > 0 {
	//	return fmt.Errorf("A total of %d errors were found. Unprocessable!", errCount)
	//}
	return errCount
}

func ValidateConfig(payload *common.Payload) {
	config := payload.Config
	_, err := common.GetFieldByName(config.Fields, "from")
	if err != nil {
		msg := "Configuration is not valid, from field is missing."
		logging.LogMessage(payload, "logs", "error", msg, "EmailConnector")
		errCount++
		//return errors.New(msg)
	}
}

func ValidateDevice(payload *common.Payload) {

	dev := payload.Device
	//fmt.Printf("####  Device #####\n%s\n", spew.Sdump(dev))
	_, err := common.GetFieldByName(dev.Fields, "to")
	if err != nil {
		msg := fmt.Sprintf("Device: %s does not have an address to delivery to.", dev.Name)
		logging.LogMessage(payload, "logs", "error", msg, "EmailConnector")
		errCount++
	}
	if strings.ToLower(dev.Active) != "true" {
		msg := fmt.Sprintf("Device: %s is not Active: %s", dev.Name, dev.Active)
		logging.LogMessage(payload, "logs", "error", msg, "EmailConnector")
		errCount++
		//return errors.New(msg)
	}
	//common.SendStatusReport("info", "checking Validation","it has never been validated",
	//	payload.DelvPayload.CallBackStatusUrl)
	fmt.Printf("ValidationDate: %s\n", dev.ValidationDate)
	if dev.ValidationDate.IsZero() {
		msg := fmt.Sprintf("Device: %s has not been validated.", dev.Name)
		logging.LogMessage(payload, "logs", "error", msg, "EmailConnector")
		errCount++
		//return errors.New(msg)
	}

	//TODO: GEt validation expiration fro config
	//TODO: Maintain list of expiring delivery methods in config
	//days:= 180
	//hours := days * 24
	//expirableMethods := []string{"FAX","EMAIL"}
	//common.SendStatusReport("info", "checking Expired","s the validation too old",payload.DelvPayload.CallBackStatusUrl)
	//
	//for _, m := range expirableMethods {
	//	if strings.ToLower(m) == dev.Method {
	//		expTime := dev.ValidatedAt.Add(4320 * time.Hour )
	//		fmt.Printf("Validation Expires at %s\n", expTime)
	//		if expTime.Before(time.Now()){
	//			msg := fmt.Sprintf("Validation for Device: %s has expired", dev.Name )
	//			fmt.Println(msg)
	//			logging.LogMessage(payload, "logs", "error", msg, "EmailConnector")
	//			errCount++
	//			//return errors.New(msg)
	//		}
	//	}
	//}

	//return nil // Valid Device
}

func ProcessEmailPayload(payload *common.Payload) error {
	name := payload.DelvPayload.Recipient.Name
	numDocs := len(payload.DelvPayload.Documents)
	//device := payload.Device
	//payload.DelvPayload.
	//fmt.Printf("\n###payload: %s\n\n", spew.Sdump(payload))
	msg := fmt.Sprintf("Create Delivery email to %s with %d docment(s)", name, numDocs)

	logging.LogMessage(payload, "logs", "info", msg, "EmailConnector")
	fmt.Printf("#### ProcessEmailPayoad ####\n")
	dDocs, err := buildDocumentSet(payload)
	if err != nil {
		msg := fmt.Sprintf("buildDocumentSet failed: %s\n", err.Error())
		logging.LogMessage(payload, "logs", "info", msg, "EmailConnector")
		return err
	}
	//fmt.Printf("#### buildDocument created: %s\n", spew.Sdump(dDocs))
	err = Send(payload, dDocs)
	if err != nil {
		fmt.Printf("Send release failed: %s\n", err.Error())
		msg := fmt.Sprintf("Send Email failed: %s\n", err.Error())
		logging.LogMessage(payload, "logs", "info", msg, "EmailConnector")

		return err
	}
	msg = fmt.Sprintf("Email was sent with %d documents", len(dDocs))
	logging.LogMessage(payload, "logs", "Success", msg, "EmailConnector")

	return nil
}

func buildDocumentSet(job *common.Payload) ([]*docMod.DeliveryDocument, error) {
	//var err error
	//fmt.Printf("\n###  BuildDocumenySet: %s\n\n\n", spew.Sdump(job))
	FileSet := []*docMod.DeliveryDocument{}
	meta := job.DelvPayload.Meta
	dob := common.GetDataByName(meta, "dob")
	patientName := common.GetDataByName(meta, "patient_name")
	names := strings.Split(patientName, ", ")
	lastName := names[0]
	//fmt.Printf("Meta: %s\n", spew.Sdump(meta))
	fmt.Printf("patient_name : %s\n", patientName)
	fmt.Printf("DOB: %s\n", dob)
	//outFileName := fmt.Sprintf("/root/tmp_images/%s.pdf", jobID)
	for _, doc := range job.DelvPayload.Documents {
		fmt.Printf("Working on first document in build set\n")
		delvDoc, err := GetDocumentFile(doc)

		if err != nil {
			msg := fmt.Sprintf("Job: %s could not retrieve document Image:[%s]", job.DelvPayload.ID, doc.ID)
			logging.LogMessage(job, "logs", "warn", msg, "email_connector")
			//TODO: If document missing on email create a place holder giving the class, description and DOS.
			continue
			//TODO: Log unable to retrieve document image for a release to an error list stored in the Delivery History
		}
		docId := doc.ID.Hex()
		tempFolder := common.GetDataByName(job.Config.Data, "temp_dir")
		if tempFolder == "" {
			tempFolder = "/data"
		}
		//tempFolder := os.Getenv("TEMP_FOLDER")
		delvDoc.FileName = fmt.Sprintf("%s/%s_%s_%s_%s_%s.%s", tempFolder, lastName, dob, delvDoc.DocClass,
			delvDoc.DateOfService, docId, delvDoc.ImageType)
		fmt.Printf("FileName: %s\n", delvDoc.FileName)
		file, err := os.OpenFile(delvDoc.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		if delvDoc.Image == nil {
			fmt.Printf("Image is nil\n")
		} else {
			fmt.Printf("Image is not nil\n")
		}
		fmt.Printf("Size of image: %d\n", len(*delvDoc.Image))
		n, err := file.Write(*delvDoc.Image)
		//if err != nil {
		//	panic(err)
		//}
		fmt.Printf("wrote %d bytes", n)
		//
		////f, _ := os.Create("./debbie.pdf")
		//w := bufio.NewWriter(f)
		//w.
		//err = ioutil.WriteFile("./debbie.pdf", *delvDoc.Image, 0644)
		if err != nil {
			fmt.Printf("Write file %s failed: %s\n", "./debbie.pdf", err.Error())
			//TODO: If document missing on email create a place holder giving the class, description and DOS.
			//TODO: Log unable to retrieve document image for a release to an error list stored in the Delivery History
			msg := fmt.Sprintf("Save Image to %s failed with %s", delvDoc.FileName, err.Error())
			logging.LogMessage(job, "logs", "error", msg, "email_connector")
		} else {
			//Since we have a file containing the image, remove the image from memory
			delvDoc.Image = nil
		}

		FileSet = append(FileSet, delvDoc)
	}
	return FileSet, nil
}

func doFunctionCall(functionCall string, data string) {
	fmt.Println("I have to do functionCall", functionCall)
	fmt.Println("for string", data)
	return
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

func makeRestCall(payload *common.Payload, method string, url string, headers map[string]string, body string) error {
	// fmt.Println("in makeRestCall", string(body))
	LogMessage(payload, "Detailed", "Info", "Calling url: "+url, payload.Config.LogURL)

	client := &http.Client{}

	if strings.Contains(body, "{{#") {
		// fmt.Println("I have handlebars")
		// remainBody := ""
		value := ""
		newValue := ""
		functionToCall := ""
		buildingNewValue := false

		// fmt.Println("split1 =", split1)
		split1 := strings.Split(body, "{{#")
		for _, item := range split1 {
			if strings.Contains(item, "#}}") {

				split2 := strings.Split(item, "#}}")

				functionToCall = split2[0]
				result, found := GetStringInBetweenTwoString(item, "#}}", "{{/"+functionToCall+"/}}")
				fmt.Println("result=", result)
				fmt.Println("found=", found)
				buildingNewValue = true
				fmt.Println("item = ", item)

				// for i, a := range split2[1] {
				// 					fmt.Println("a", i, a)
				// 				}
				rmnStr := split2[1]
				remainBody := strings.Split(rmnStr, "{")

				// test1 := strings.IndexByte(rmnStr, byte("{{/"+functionToCall+"/}}"))
				// 				fmt.Println("test1==========================================================")
				// 				fmt.Println(test1)
				// fmt.Println("functionToCall=", functionToCall)
				// fmt.Println("reaminBody=", remainBody)

				for _, remain := range remainBody {
					if buildingNewValue == true {
						newValue = newValue + " " + remain
						if strings.HasPrefix(remain, "{/"+functionToCall+"/}}") {

							fmt.Println("found end of function")
							buildingNewValue = false
							doFunctionCall(functionToCall, newValue)
						}
					} else {
						value = value + " " + remain
					}
				}

				//
				// 				// fmt.Println("I have item=", item)
				// 				fmt.Println("I see functionCall", functionToCall)
				// 				fmt.Println("remainBody=", remainBody)
				// 				if strings.HasPrefix(item, "#") {
				// 					fmt.Println("I have a begining escape seq", item)
				// 					buildingNewValue = true
				// 					// value = value + " {{" + item
				// 				}
				// 				// if strings.ToLower(strings.Split(rflctVal[0], ".")[0]) == "response" && actionResponse != nil {
				// 				// 					value = value + reflectActionResponseFieldValue(actionResponse, rflctVal[0])
				// 				// 					for key, val := range rflctVal {
				// 				// 						if key != 0 {
				// 				// 							value = value + val
				// 				// 						}
				// 				// 					}
				// 				// 				} else {
				// 				// 					value = value + reflectValue(alert, rflctVal[0])
				// 				// 					// fmt.Println("I have value:", value)
				// 				// 					for key, val := range rflctVal {
				// 				// 						if key != 0 {
				// 				// 							value = value + val
				// 				// 						}
				// 				// 					}
				// 				// 				}
				// 			} else if buildingNewValue == true {
				// 				newValue = newValue + " " + item
				// 			} else {
				// 				value = value + " " + item
			} else {
				value = value + " " + item
			}
		}

		// fmt.Println("field value =", value)
		fmt.Println("value=", value)
	}

	return nil

	var bstr []byte

	// if len(headers) > 0 {
	// 		for k, h := range headers {
	// 			fmt.Println("h=", h)
	// 			fmt.Println("k=", k)
	// 			if strings.ToLower(k) == "content-type" {
	// 				if strings.Contains(h, "json") {
	// 					fmt.Println("I see I have json")
	// 					tbstr, err := json.Marshal(body)
	// 					if err != nil {
	// 						return err
	// 					}
	// 					bstr = tbstr
	// 				}
	// 			} else {
	// 				bstr = []byte(body)
	// 			}
	// 		}
	// 	} else {
	bstr = []byte(body)
	// }

	fmt.Println(string(bstr))

	// return nil

	req, _ := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(bstr))
	if len(headers) > 0 {
		for k, h := range headers {
			// fmt.Println("h=", h)
			// fmt.Println("k=", k)
			req.Header.Set(k, h)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		LogMessage(payload, "Detailed", "Error", "Error when calling: "+url+" | received status code: "+strconv.Itoa(res.StatusCode), payload.Config.LogURL)
		return err
	}
	defer res.Body.Close()
	fmt.Println("res", res.StatusCode)
	if res.StatusCode != 200 {
		LogMessage(payload, "Basic", "Error", "Received Error Code "+strconv.Itoa(res.StatusCode)+" when calling url:"+url, payload.Config.LogURL)
		return nil
	}

	LogMessage(payload, "Basic", "Success", "Successfully called url:"+url+" | received status code: "+strconv.Itoa(res.StatusCode), payload.Config.LogURL)
	return nil
}

//CreateRestCall
func CreateDelivery(payload *common.Payload) error {
	//fmt.Printf("in CreateDelivery: %s\n", spew.Sdump(payload))
	// spew.Dump(payload.Action.Fields)
	//LogMessage(payload, "Basic", "Success", "Delivery received.  Processing....", payload.Config.LogURL)
	//toField, _ := m.GetDeviceField(payload., "to")
	return nil
}

// //toField, _ := m.GetDeviceField(payload., "to")

// var body string
// var headers string
// var url string
// method := "POST"

// for _, c := range payload.Action.Fields {
// 	switch {
// 	case c.Name == "body":
// 		val := len(strings.TrimSpace(c.Value))
// 		if val > 0 {
// 			body = strings.TrimSpace(c.Value)
// 			LogMessage(payload, "Detailed", "Info", "Found "+c.Name+" field: "+body, payload.Config.LogURL)
// 		} else {
// 			LogMessage(payload, "Detailed", "Info", c.Name+" field is empty", payload.Config.LogURL)
// 		}

// 	case c.Name == "headers":
// 		val := len(strings.TrimSpace(c.Value))
// 		if val > 0 {
// 			headers = strings.TrimSpace(c.Value)
// 			LogMessage(payload, "Detailed", "Info", "Found "+c.Name+" field: "+headers, payload.Config.LogURL)
// 		} else if len(strings.TrimSpace(c.Display_value)) > 0 {
// 			headers = strings.TrimSpace(c.Display_value)
// 			LogMessage(payload, "Detailed", "Info", "Found "+c.Name+" display value: "+headers, payload.Config.LogURL)
// 		} else {
// 			LogMessage(payload, "Detailed", "Info", c.Name+" field is empty", payload.Config.LogURL)
// 		}
// 	case c.Name == "method":
// 		val := len(strings.TrimSpace(c.Value))
// 		if val > 0 {
// 			method = strings.TrimSpace(c.Value)
// 			LogMessage(payload, "Detailed", "Info", "Found "+c.Name+" field: "+method, payload.Config.LogURL)
// 		} else {
// 			LogMessage(payload, "Detailed", "Info", c.Name+" field is empty", payload.Config.LogURL)
// 		}
// 	case c.Name == "url":
// 		val := len(strings.TrimSpace(c.Value))
// 		if val > 0 {
// 			url = strings.TrimSpace(c.Value)
// 			LogMessage(payload, "Detailed", "Info", "Found "+c.Name+" field: "+url, payload.Config.LogURL)
// 		} else {
// 			LogMessage(payload, "Detailed", "Info", c.Name+" field is empty", payload.Config.LogURL)
// 		}
// 	}
// }

// fmt.Println("I have body:", body)
// fmt.Println("I have headers:", headers)
// fmt.Println("I have method:", method)
// fmt.Println("I have url:", url)

// var bBody []byte
//
// 	if len(body) > 0 {
// 		bBody, err := json.Marshal(body)
// 		if err != nil {
// 			fmt.Println("error:", err.Error())
// 			// return bstr, err
// 		}
// 		fmt.Println("json body:", bBody)
// 	}

// var headerMap map[string]string
// // We'll store the error (if any) so we can return it if necessary
// // fmt.Println("before header marshal")
// if len(headers) > 0 {
// 	err := json.Unmarshal([]byte(headers), &headerMap)
// 	if err != nil {
// 		log.Println("Error Unmarshaling headers")
// 		return err
// 	}
// }

// // fmt.Println("after header marshal")
// // fmt.Println(headerMap)

// err := makeRestCall(payload, method, url, headerMap, body)
// if err != nil {
// 	fmt.Println("error:", err.Error())
// 	// LogMessage(payload, "Detailed", "Error", "Error when calling url: "+ url+" | Error: "+err.Error(), payload.Config.Core_log_url)
// }
// // LogMessage(payload, "Basic", "Success", "Incident successfully created. number: "+ticket+", sys_id: "+sys_id, payload.Config.Core_log_url)
// if len(payload.Action.Subactions) > 0 {
// 	var respFlds ActionResponseFields
// 	CallProcessSubactions(payload, &respFlds)
// }

// return nil
//}

// func CallProcessSubactions(payload *common.Payload, actionResponseFields *ActionResponseFields) {
// 	LogMessage(payload, "Detailed", "Info", "Sending request to core to process subactions", payload.Config.Core_log_url)
// 	url := payload.Config.LogURL + "/action/" + payload.Action.Id.Hex() + "/alert/" + payload.Event_id.Hex() + "/subscription/" + payload.Subscription_id.Hex()
// 	fmt.Println("subactions ur=", url)

// 	jsonPayload, err := json.Marshal(actionResponseFields)
// 	if err != nil {
// 		log.Println("Error marshaling json in CallProcessSubactions", err.Error())
// 		return
// 	}

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		fmt.Println("Error when calling process subactions", err.Error())
// 		return
// 	}

// 	defer resp.Body.Close()
// 	return
// }
