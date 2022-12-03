package document

type Document struct {
	DownloadURL string
	DocumentID  string
	CategoryID  string
	PatientID   string
}

func (d *Document) FileName() string {
	return d.PatientID
}
