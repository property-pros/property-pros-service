package documents

import "github.com/vireocloud/property-pros-service/interfaces"

type DocumentContent struct {
	interfaces.IDocumentContent
	DocSource                 []byte
	DocTemplate               []byte
	SerializedTemplateModel   []byte
	DeserializedTemplateModel interface{}
	DocContent                []byte
	UserId                    string
}

func (d *DocumentContent) GetDocContent() []byte {
	return d.DocContent
}
