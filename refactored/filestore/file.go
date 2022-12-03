package filestore

import (
	"fmt"
	"io"
	"os"

	"github.com/bradleyshawkins/go-clean-architecture/refactored/document"
)

type TempFiler interface {
	NewTempFile(doc *document.Document) (io.ReadWriteCloser, error)
}

type TempFile struct {
	file *os.File
}

func (t *TempFile) Read(p []byte) (n int, err error) {
	return t.file.Read(p)
}

func (t *TempFile) Write(p []byte) (n int, err error) {
	return t.file.Write(p)
}

func (t *TempFile) Close() error {
	err := t.file.Close()
	if err != nil {
		return err
	}

	fmt.Println("Removing temp file")
	//err = os.Remove(t.file.Name())
	//if err != nil {
	//	return err
	//}

	return nil
}

type OS struct {
	directory string
}

func NewOS(directory string) *OS {
	return &OS{directory: directory + "/temp-documents"}
}

func (o *OS) NewTempFile(doc *document.Document) (io.ReadWriteCloser, error) {
	err := o.ensureDirectoryExists()
	if err != nil {
		return nil, err
	}

	pattern := determinePattern(doc)

	fmt.Println("Creating temp file at", o.directory)
	tmpFile, err := os.CreateTemp(o.directory, pattern)
	if err != nil {
		return nil, err
	}

	fmt.Println("Created temp file named:", tmpFile.Name())
	return &TempFile{file: tmpFile}, nil
}

func (o *OS) ensureDirectoryExists() error {
	err := os.MkdirAll(o.directory, 0777)
	if err != nil {
		return err
	}
	return nil
}

// determinePattern could be improved to figure out file extensions instead of saving everything as jpg
func determinePattern(doc *document.Document) string {
	return fmt.Sprintf("%s-*.jpg", doc.PatientID)
}
