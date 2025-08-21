package handler

import (
	"fmt"
	"net/http"

	"github.com/Dsouza10082/documentOCRProcessor/model"
)

func PythonTextExtractorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("📋 File extraction started\n")
	fp := model.NewFileProcessor()
	filePaths, _ := fp.ReadDirectory()
	if len(filePaths) == 0 {
		fmt.Printf("📭 No files found in directory '%s'\n", fp.ToProcessDir)
	}
	fmt.Printf("📋 Found %d file(s) to process:\n", len(filePaths))
	for i, path := range filePaths {
		fmt.Printf("   %d. %s\n", i+1, path)
	}
	fp.ProcessFiles(filePaths)
}