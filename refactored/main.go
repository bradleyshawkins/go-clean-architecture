package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bradleyshawkins/go-clean-architecture/refactored/document"
	"github.com/bradleyshawkins/go-clean-architecture/refactored/download"
	"github.com/bradleyshawkins/go-clean-architecture/refactored/filestore"
	"github.com/bradleyshawkins/go-clean-architecture/refactored/handlers"
	"github.com/bradleyshawkins/go-clean-architecture/refactored/integrations"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	path := filepath.Dir(os.Args[0])
	fmt.Println(path)

	integration := &integrations.Implementation{}

	osFileStore := filestore.NewOS(path)

	fileDownloader := download.NewFileDownloader(client, osFileStore)

	documentImporter := document.NewImporter(integration, fileDownloader)

	fileDownloadHandler := handlers.NewFileDownloadHandler(documentImporter)

	handlers.Start(fileDownloadHandler)
}
