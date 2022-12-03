package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradleyshawkins/go-clean-architecture/refactored/document"
)

type FileDownloadHandler struct {
	importer *document.Importer
}

func NewFileDownloadHandler(importer *document.Importer) *FileDownloadHandler {
	return &FileDownloadHandler{importer: importer}
}

type downloadFileRequest struct {
	DownloadURL string `json:"downloadURL"`
	DocumentID  string `json:"documentID"`
	CategoryID  string `json:"categoryID"`
	PatientID   string `json:"patientID"`
}

func (f *FileDownloadHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	var req downloadFileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Printf("Received request to download file from %s...\n", req.DownloadURL)

	doc := &document.Document{
		DownloadURL: req.DownloadURL,
		DocumentID:  req.DocumentID,
		CategoryID:  req.CategoryID,
		PatientID:   req.PatientID,
	}

	err = f.importer.ImportDocument(r.Context(), doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
