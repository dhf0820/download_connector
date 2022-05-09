package document

import (
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergeToCombined(inFiles []string, outFileName string) error {

	err := api.MergeAppendFile(inFiles, outFileName, nil) //Keep appending new documents to combined file
	for _, fn := range inFiles {
		fmt.Printf("Removing file: [%s]\n", fn)
		os.Remove(fn)
	}
	return err
}
