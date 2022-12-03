package integrations

import (
	"context"
	"fmt"
	"io"

	"github.com/bradleyshawkins/go-clean-architecture/refactored/document"
)

type Implementation struct {
}

func (i *Implementation) ImportDocument(ctx context.Context, doc *document.Document, file io.Reader) error {
	fmt.Println("I'm importing a document!")
	return nil
}
