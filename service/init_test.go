package service

import (
	//coreMod "gitlab.com/dhf0820/ids_model/coreMod"
	"os"
)

//https://stackoverflow.com/questions/42102496/testing-a-grpc-service
const bufSize = 1024 * 1024

const chunkSize = uint32(1 << 14)

func InitTest() {
	//fmt.Printf("InitTest\n")
	os.Setenv("CONFIG_ADDRESS", "http://localhost:29900/api/v1/")
	os.Setenv("SERVICE_NAME", "delivery")

	os.Setenv("SERVICE_VERSION", "local")
	os.Setenv("COMPANY", "demo")
	//cfg, err := Initialize("local", "demo")
	//if err != nil {
	//	os.Exit(2)
	//}
	////	fmt.Printf("cfg: %s\n",spew.Sdump(cfg))
	//return cfg
}

//func init() {
//
//	bufLis = bufconn.Listen(bufSize)
//	s := grpc.NewServer()
//	delPB.RegisterDeliveryServiceServer(s, &DeliveryServiceServer{})
//	go func() {
//		if err := s.Serve(bufLis); err != nil {
//			log.Fatalf("Server exited with error: %v", err)
//		}
//	}()
//
//}
//
//func bufDialer(context.Context, string) (net.Conn, error) {
//	return bufLis.Dial()
//}

/*func SetupDomainRelease(t *testing.T, create bool) *domain.Release {
	//settings.SetDbName("test_documents")
	data := sample.NewDomainDocument(1)
	//data.ID = primitive.NilObjectID
	// if create {
	// 	fmt.Printf("Creating new document: %s\n", spew.Sdump(data))
	// 	doc, err = domain.AddDocument(data)
	// 	if err != nil {
	// 		t.Fatalf("Error setupDomainDocument Creating: %v", err)
	// 	}
	// } else {
	// 	doc = data
	// }
	return data
}

func SetupPbCreateRelease(t *testing.T) *pb.CreateRelease {
	//settings.SetDbName("test_documents")
	data := sample.NewDocument(1)
	//imageFile := sample.ImageFileName
	//_, err := data.FromDocumentPB(pbDoc)
	// if err != nil {
	// 	err := fmt.Errorf("FromDocumentPB failed: %v", err)
	// 	t.Fatal(err)
	// }
	return data
}*/
