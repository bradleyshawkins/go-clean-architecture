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
	FileName    string `json:"fileName"`
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
		FileName:    req.FileName,
		DocumentID:  req.DocumentID,
		PatientID:   req.PatientID,
	}

	err = f.importer.ImportDocument(r.Context(), doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
