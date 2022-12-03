package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bradleyshawkins/go-clean-architecture/existing/handlers"
	"github.com/bradleyshawkins/go-clean-architecture/existing/integrations"
)

func main() {
	impl := &integrations.Implementation{}
	path := filepath.Dir(os.Args[0])
	fmt.Println(path)
	fileDownloader := handlers.NewFileDownloader(impl, path, &http.Client{})

	handlers.Start(fileDownloader)
}
