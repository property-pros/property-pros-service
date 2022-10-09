package agreements

import (
	"log"

	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementModel struct {
	*interop.NotePurchaseAgreement
	*common.BaseModel[interop.NotePurchaseAgreement]
	documentContentService       interfaces.IDocumentContentService
	notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway
	userService                  interfaces.IUsersService
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) Save() (interfaces.IAgreementModel, error) {
	// docChannel := make(chan interfaces.IDocumentContent)
	// userChannel := make(chan interfaces.IUserModel)

	if exists, err := notePurchaseAgreement.DoesExist(); err == nil && !exists {
		// go notePurchaseAgreement.GenerateDocument(docChannel)
	} else {
		// close(docChannel)
		return interfaces.IAgreementModel(notePurchaseAgreement), err
	}

	go notePurchaseAgreement.SaveUser(nil)
	go notePurchaseAgreement.SaveNotePurchaseAgreement(nil)

	return notePurchaseAgreement, nil
}

// func (notePurchaseAgreement *NotePurchaseAgreementModel) GetContext() context.Context {
// 	return notePurchaseAgreement.BaseModel.GetContext()
// }

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

func (notePurchaseAgreement *NotePurchaseAgreementModel) SaveUser(resultChannel chan<- *interop.User) (*interop.User, error) {

	result, err := notePurchaseAgreement.userService.SaveUser(notePurchaseAgreement.Context, notePurchaseAgreement.GetUserPayload())
	if err != nil {
		close(resultChannel)
		return nil, err
	}

	if resultChannel != nil {
		resultChannel <- result
	}

	return result, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) SaveNotePurchaseAgreement(docChannel <-chan interfaces.IDocumentContent) (interfaces.IAgreementModel, error) {

	result, err := notePurchaseAgreement.notePurchaseAgreementGateway.SaveNotePurchaseAgreement(notePurchaseAgreement.GetContext(), notePurchaseAgreement)

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

func (agreement *NotePurchaseAgreementModel) GetPassword() string {
	return ""
}

func (agreement *NotePurchaseAgreementModel) GetUserPayload() *interop.User {
	return &interop.User{
		EmailAddress: agreement.GetEmailAddress(),
		Password:     agreement.GetPassword(),
	}
}

func NewNotePurchaseAgreementModel(
	documentContentService interfaces.IDocumentContentService,
	notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway,
	userService interfaces.IUsersService,
) *NotePurchaseAgreementModel {
	log.Printf("building note purchase agreement model; documentservice: %+#v\n\n", documentContentService)
	return &NotePurchaseAgreementModel{
		documentContentService:       documentContentService,
		notePurchaseAgreementGateway: notePurchaseAgreementGateway,
		userService:                  userService,
	}
}
