package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Start(f *FileDownloadHandler) {
	m := chi.NewMux()

	m.Post("/download", f.DownloadFile)

	fmt.Println("Starting router...")
	err := http.ListenAndServe(":8080", m)
	if err != nil {
		fmt.Printf("Error received: %v\n", err)
	}
}
