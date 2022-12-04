package document

type Document struct {
	DownloadURL string
	FileName    string
	DocumentID  string
	PatientID   string
}

func (d *Document) FileExtension() string {
	return ".png"
}
