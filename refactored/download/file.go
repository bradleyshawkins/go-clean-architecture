package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bradleyshawkins/go-clean-architecture/refactored/document"
	"github.com/bradleyshawkins/go-clean-architecture/refactored/filestore"
)

type File struct {
	client    *http.Client
	tempFiler filestore.TempFiler
}

func NewFileDownloader(client *http.Client, tempFiler filestore.TempFiler) *File {
	return &File{
		client:    client,
		tempFiler: tempFiler,
	}
}

// DownloadToTempFile downloads the document provided to a temp file and returns an io.ReadCloser so it can be read
// and cleaned up
func (f *File) DownloadToTempFile(ctx context.Context, doc *document.Document) (io.ReadCloser, error) {
	u, err := url.Parse(doc.DownloadURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url provided. %w", err)
	}

	fmt.Println("Creating new temp file...")
	tempFile, err := f.tempFiler.NewTempFile(doc)
	if err != nil {
		return nil, err
	}

	fmt.Println("Downloading to temp file")
	err = f.downloadTo(u, tempFile)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func (f *File) downloadTo(u *url.URL, tempFile io.ReadWriteCloser) error {

	req, err := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("unable to create request. %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request to download file. %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code received. %w StatusCode: %d", err, resp.StatusCode)
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to copy payload to writer. %w", err)
	}

	return nil
}
