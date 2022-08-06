package agreements

import (
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementModel struct {
	*interop.BaseModel[interop.NotePurchaseAgreement]
	documentContentService       interfaces.IDocumentContentService
	notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway
	userService                  interfaces.IUserService
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) Save() {
	docChannel := make(chan interfaces.IDocumentContent)
	userChannel := make(chan interfaces.IUserModel)

	if exists, err := notePurchaseAgreement.DoesExist(); err == nil && exists {
		go notePurchaseAgreement.GenerateDocument(docChannel)
	} else {
		close(docChannel)
	}

	go notePurchaseAgreement.SaveUser(userChannel)
	go notePurchaseAgreement.SaveNotePurchaseAgreement(docChannel)
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) DoesExist() (bool, error) {
	return false, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) GenerateDocument(resultChannel chan<- interfaces.IDocumentContent) (interfaces.IDocumentContent, error) {

	result, err := notePurchaseAgreement.documentContentService.BuildNotePurchaseAgreement(notePurchaseAgreement.Context, notePurchaseAgreement.Payload)

	if err != nil {
		close(resultChannel)
		return nil, err
	}

	resultChannel <- result

	return result, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) SaveUser(resultChannel chan<- interfaces.IUserModel) (interfaces.IUserModel, error) {

	result, err := notePurchaseAgreement.userService.SaveUser(notePurchaseAgreement.Context, notePurchaseAgreement.GetUser())

	if err != nil {
		close(resultChannel)
		return nil, err
	}

	resultChannel <- result

	return result, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) SaveNotePurchaseAgreement(docChannel <-chan interfaces.IDocumentContent) (interfaces.IAgreementModel, error) {

	result, err := notePurchaseAgreement.notePurchaseAgreementGateway.SaveNotePurchaseAgreement(notePurchaseAgreement.Context, notePurchaseAgreement)

	if err != nil {
		return nil, err
	}

	// doc := <-docChannel

	// doc.SaveDocument()

	return result, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) GetDocumentContent(resultChannel chan<- []byte) ([]byte, error) {
	doc, err := notePurchaseAgreement.documentContentService.BuildNotePurchaseAgreement(notePurchaseAgreement.Context, notePurchaseAgreement.Payload)

	if err != nil {
		if resultChannel != nil {
			close(resultChannel)
		}

		return nil, err
	}

	result := doc.GetDocContent()

	if resultChannel != nil {
		resultChannel <- result
	}

	return result, nil
}

func (agreement *NotePurchaseAgreementModel) GetUser() interfaces.IUserModel {
	return agreement.Payload
}

func NewNotePurchaseAgreementModel(documentContentService interfaces.IDocumentContentService) *NotePurchaseAgreementModel {
	return &NotePurchaseAgreementModel{
		documentContentService: documentContentService,
	}
}
