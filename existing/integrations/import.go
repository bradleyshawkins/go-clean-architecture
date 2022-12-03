package integrations

import (
	"fmt"
	"github.com/bradleyshawkins/go-clean-architecture/existing/handlers"
	"io"
)

type Implementation struct {
}

func (i *Implementation) ImportDocument(doc *handlers.DownloadDocument, file io.Reader) error {
	fmt.Println("I'm importing a document!")
	return nil
}
