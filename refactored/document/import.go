package document

import (
	"context"
	"fmt"
	"io"
)

type documentImporter interface {
	ImportDocument(ctx context.Context, document *Document, r io.Reader) error
}

type downloader interface {
	DownloadToTempFile(ctx context.Context, document *Document) (io.ReadCloser, error)
}

type Importer struct {
	documentImporter documentImporter
	fileDownloader   downloader
}

func NewImporter(documentImporter documentImporter, downloader downloader) *Importer {
	return &Importer{
		documentImporter: documentImporter,
		fileDownloader:   downloader,
	}
}

func (i *Importer) ImportDocument(ctx context.Context, doc *Document) error {
	fmt.Println("Downloading to temp file...")
	tempFile, err := i.fileDownloader.DownloadToTempFile(ctx, doc)
	if err != nil {
		return err
	}

	defer tempFile.Close()

	fmt.Println("Importing document...")
	err = i.documentImporter.ImportDocument(ctx, doc, tempFile)
	if err != nil {
		return err
	}

	return nil

}
