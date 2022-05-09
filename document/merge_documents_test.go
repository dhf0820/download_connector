package document

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	//docMod "gitlab.com/dhf0820/ids_model/document"

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
