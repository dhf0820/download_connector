package document

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gitlab.com/dhf0820/ids_model/common"
	//docMod "gitlab.com/dhf0820/ids_model/document"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
	//deliv "gitlab.com/dhf0820/delivery_service/connect_deliver"
	//pb "gitlab.com/dhf0820/delivery_service/protobufs/delPB"
	"testing"
)

func TestGetDocument(t *testing.T) {
	t.Parallel()
	fmt.Printf("----TestGetDocument\n")
	InitTest()
	Convey("GetDocument via Restful", t, func() {
		docMeta := sampleDocMeta()
		doc, err := GetDocumentWithImage(docMeta)
		So(err, ShouldBeNil)
		So(doc, ShouldNotBeNil)
		fmt.Printf("ReceivedDocument: %s\n", spew.Sdump(doc))
	})
}

func sampleDocMeta() *common.DocumentMeta {
	docID, _ := primitive.ObjectIDFromHex("61b6c89bd8cc9541dde09f1b")
	docMeta := common.DocumentMeta{}
	docMeta.ID = docID
	docMeta.DocumentURL = "http://localhost:29912/api/v1/document/61b6c89bd8cc9541dde09f1b?image=include"
	docMeta.Description = "X-Ray Left Knee"
	docMeta.Meta = []*common.KVData{}
	cls := common.KVData{
		Name:  "doc_class",
		Value: "Radiology",
	}
	dos := common.KVData{
		Name:  "date_of_service",
		Value: time.Now().String(),
	}
	docMeta.Meta = append(docMeta.Meta, &cls)
	docMeta.Meta = append(docMeta.Meta, &dos)
	docMeta.Status = &common.Status{
		State:      "",
		StatusTime: time.Now(),
		Comment:    "ok",
	}

	return &docMeta

}
