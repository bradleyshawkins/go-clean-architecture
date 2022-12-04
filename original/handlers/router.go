package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Start(downloader *FileDownloader) {
	m := chi.NewMux()

	m.Post("/download", downloader.DownloadFileHandler)

	fmt.Println("Starting router...")
	err := http.ListenAndServe(":8080", m)
	if err != nil {
		fmt.Printf("Error received: %v\n", err)
	}
}
