package document

import (
	"encoding/json"
	"fmt"

	//"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	common "gitlab.com/dhf0820/ids_model/common"
	docMod "gitlab.com/dhf0820/ids_model/document"
)

func GetDocumentFile(docMeta *common.DocumentMeta) (*docMod.DeliveryDocument, error) {
	fmt.Printf("\n\n###  Building DeliveryDocument:16 ###\n\n")
	documentResp := new(docMod.DocumentResponse)
	docUrl := fmt.Sprintf("%s?image=include&format=pdf", docMeta.DocumentURL)
	fmt.Printf("GetDocumentUrl: %s\n", docUrl)
	req, err := http.NewRequest("GET", docUrl, nil)
	if err != nil {
		log.Errorf("Get Document %s failed: %s", docUrl, err.Error())
		return nil, err
	}

	//fmt.Printf("Access Authorization: %s\n", endPoint.Access)
	//req.Header.Set("AUTHORIZATION", "37")
	//req.Header.Set("facility", "demo")
	client := &http.Client{Timeout: time.Second * 10}
	//fmt.Printf("Calling for document: %s\n\n", spew.Sdump(docMeta))
	resp, err := client.Do(req)
	if err != nil {
		//TODO: Handle Error from getting document
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	fmt.Printf("reading the body doc:35\n")
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("400|document read body: %s", err.Error())
		// log.Errorf("%s", err.Error())
		// return nil, err
	}
	fmt.Printf("documentResponse:\n%s\n", string(b))
	///mt.Printf("Unmarshaling body\n")
	if err := json.Unmarshal(b, &documentResp); err != nil {
		err = fmt.Errorf("400|Error marshaling JSON: %s", err.Error())
		log.Errorf("DocumentResponse:50 --  error: %s", err)
		return nil, err
	}
	// fmt.Printf("Reaady to unmarshal the document Requested")
	// decode := json.NewDecoder(resp.Body).Decode(&documentResp)
	// if decode != nil {
	// 	err := fmt.Errorf("error Decoding GetDocumentFile:37 %s", decode.Error())
	// 	//TODO: Handle decode document error properly
	// 	log.Fatal(err)
	// }
	fmt.Printf("\n\n###  Building DeliveryDocument ###\n\n")
	doc := documentResp.Document
	//fmt.Printf("documentResp:58 %s\n", spew.Sdump(documentResp.Document))
	dd := docMod.DeliveryDocument{}

	dd.Image = documentResp.Image
	dd.Meta = docMeta.Meta
	dd.Description = common.GetDataByName(docMeta.Meta, "description")
	dd.DocClass = common.GetDataByName(docMeta.Meta, "doc_class")
	dd.DateOfService = common.GetDataByName(docMeta.Meta, "date_of_service")
	//dd.DateOfService = common.GetDataByName(dd.Meta, "date_of_service")
	dd.ImageType = doc.ImageType
	dd.ImageRepository = doc.ImageRepository
	dd.DocumentURL = docMeta.DocumentURL

	dd.ImageURL = fmt.Sprintf("%s?image=only", docMeta.DocumentURL)
	//fmt.Printf("\n\n#### DeliveryDoc:%s\n\n### docMeta: %s\n", spew.Sdump(dd), spew.Sdump(docMeta))
	//dd.ImageURL = docMeta.DocumentURL + //This gets teh combined image and document.
	//dd.DocumentURL = "http://localhost:29912/api/v1/document/61b6c89bd8cc9541dde09f1b?image=none"

	return &dd, nil

}

func GetDocumentWithImage(docMeta *common.DocumentMeta) (*docMod.DeliveryDocument, error) {
	documentResp := docMod.DocumentResponse{}
	docUrl := docMeta.DocumentURL

	req, err := http.NewRequest("GET", docUrl, nil)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Access Authorization: %s\n", endPoint.Access)
	//req.Header.Set("AUTHORIZATION", "37")
	//req.Header.Set("facility", "demo")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		//TODO: Handle Error from getting document
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	decode := json.NewDecoder(resp.Body).Decode(&documentResp)
	if decode != nil {
		err := fmt.Errorf("error Decoding GetDocumentWithImage: 77 document: %s", decode.Error())
		//TODO: Handle decoe document error properly
		log.Fatal(err)
	}

	dd := docMod.DeliveryDocument{}
	dd.Description = docMeta.Description
	dd.Image = documentResp.Image
	dd.DocClass = common.GetDataByName(docMeta.Meta, "doc_class")
	dd.DateOfService = common.GetDataByName(docMeta.Meta, "date_of_service")
	dd.ImageType = documentResp.Document.ImageType
	dd.ImageRepository = documentResp.Document.ImageRepository
	dd.ImageURL = docMeta.DocumentURL //This gets teh combined image and document.
	dd.DocumentURL = "http://localhost:29912/api/v1/document/61b6c89bd8cc9541dde09f1b?image=none"
	//dd.FileName = fmt.Sprintf("doc-%s-%s.dd", docMeta.ID.Hex, dd.DocClass)
	//ioutil.WriteFile(dd.FileName, *image,0444)
	//fmt.Printf("dd: %s\n", spew.Sdump(dd))
	dd.Image = documentResp.Image
	return &dd, nil

}
