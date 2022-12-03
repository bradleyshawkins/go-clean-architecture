package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type DocumentImporter interface {
	ImportDocument(doc *DownloadDocument, file io.Reader) error
}

type FileDownloader struct {
	docImporter  DocumentImporter
	downloadPath string
	client       *http.Client
}

func NewFileDownloader(docImporter DocumentImporter, downloadPath string, client *http.Client) *FileDownloader {
	return &FileDownloader{
		docImporter:  docImporter,
		downloadPath: downloadPath,
		client:       client,
	}
}

type DownloadDocument struct {
	DownloadURL string `json:"downloadURL"`
	DocumentID  string `json:"documentID"`
	CategoryID  string `json:"categoryID"`
	PatientID   string `json:"patientID"`
}

func (f *FileDownloader) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	var req DownloadDocument
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Received request to download file from:", req.DownloadURL)

	err = f.ImportDocument(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (f *FileDownloader) ImportDocument(doc *DownloadDocument) error {

	if !f.canInsertDocument() {
		return errors.New("cannot insert document")
	}

	fmt.Println("Creating temp file...")

	// get absolute path for the application
	aPath, err := filepath.Abs(f.downloadPath)
	if err != nil {
		return fmt.Errorf("unable to get filepath. %w", err)
	}

	parentDir := filepath.Dir(aPath)

	destDir := filepath.Join(parentDir, "document")

	err = os.MkdirAll(destDir, 0777)
	if err != nil {
		return fmt.Errorf("unable to create directories. %w", err)
	}

	pattern := fmt.Sprintf("%s-*.jpg", doc.PatientID)
	tmpFile, err := os.CreateTemp(destDir, pattern)
	if err != nil {
		return fmt.Errorf("unable to create tempFile. %w", err)
	}

	// ensure document is removed so nothing is left behind
	defer func() error {
		fmt.Println("Cleaning up temp file...")

		if err := tmpFile.Close(); err != nil {
			fmt.Println(fmt.Errorf("unable to close temp file. %w", err))
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			return fmt.Errorf("unable to remove temp file. %w", err)
		}
		return nil
	}()

	u, err := f.getDocumentURL(doc.DownloadURL)
	if err != nil {
		return err
	}

	fmt.Println("Downloading file...")

	err = f.downloadFile(f.client, u, tmpFile)
	if err != nil {
		return err
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("unable to seek to beginning of temp file. %w", err)
	}

	fmt.Println("Importing document...")

	err = f.docImporter.ImportDocument(doc, tmpFile)
	if err != nil {
		return fmt.Errorf("unable to import document. %w", err)
	}

	return nil
}

func (f *FileDownloader) canInsertDocument() bool {
	return true
}

func (f *FileDownloader) getDocumentURL(downloadURL string) (*url.URL, error) {

	u, err := url.Parse(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid download url provided. %w", err)
	}
	return u, nil
}

func (f *FileDownloader) downloadFile(c *http.Client, u *url.URL, w io.Writer) error {

	req, err := http.NewRequest(http.MethodGet, u.String(), http.NoBody)

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request to download file. %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code received. %w StatusCode: %d", err, resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to copy payload to writer. %w", err)
	}

	return nil
}
