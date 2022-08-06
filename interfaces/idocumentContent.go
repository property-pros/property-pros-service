package interfaces

type IDocumentContent interface {
	GetDocSource() []byte
	GetDocContent() []byte
}
