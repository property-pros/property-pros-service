package agreements

import (
	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/users"
)

type NotePurchaseAgreementModel struct {
	*common.BaseModel[interop.NotePurchaseAgreement]
	Id             string          `json:"id,omitempty"`
	FirstName      string          `json:"first_name,omitempty"`
	LastName       string          `json:"last_name,omitempty"`
	DateOfBirth    string          `json:"date_of_birth,omitempty"`
	HomeAddress    string          `json:"home_address,omitempty"`
	User           users.UserModel `json:"user,omitempty" gorm:"foreignKey:UserId;references:Id"`
	PhoneNumber    string          `json:"phone_number,omitempty"`
	SocialSecurity string          `json:"social_security,omitempty"`
	FundsCommitted uint64          `json:"funds_committed,omitempty"`
	FileContent    []byte          `json:"file_content,omitempty"`
	CreatedOn      string          `json:"created_on,omitempty"`

	UserId                       string `json:"userid" bson:"userid"`
	documentContentService       interfaces.IDocumentContentService
	notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway
	userService                  interfaces.IUsersService
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) Save() (interfaces.IAgreementModel, error) {
	// docChannel := make(chan interfaces.IDocumentContent)
	// userChannel := make(chan interfaces.IUserModel)

	// if exists, err := notePurchaseAgreement.DoesExist(); err == nil && !exists {
	// 	go notePurchaseAgreement.GenerateDocument(docChannel)
	// } else {
	// 	close(docChannel)
	// }
	userIdString := uuid.New().String()
	notePurchaseAgreement.User.Id = userIdString
	notePurchaseAgreement.UserId = userIdString
	go notePurchaseAgreement.SaveUser(nil)
	// go notePurchaseAgreement.SaveNotePurchaseAgreement(nil)

	return notePurchaseAgreement, nil
}

func (notePurchaseAgreement *NotePurchaseAgreementModel) LoadDocument() (interfaces.IAgreementModel, error) {
	result := *notePurchaseAgreement

	docContent, err := result.GenerateDocument(nil)

	if err != nil {
		return nil, err
	}

	result.GetPayload().FileContent = docContent.GetDocContent()

	return &result, nil
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

	if resultChannel != nil {
		resultChannel <- result
	}

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
	return agreement.GetPayload().User
}

func (agreement *NotePurchaseAgreementModel) GetId() string {
	return agreement.GetPayload().Id
}

func NewNotePurchaseAgreementModel(
	documentContentService interfaces.IDocumentContentService,
	notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway,
	userService interfaces.IUsersService,
) *NotePurchaseAgreementModel {

	return &NotePurchaseAgreementModel{
		documentContentService:       documentContentService,
		notePurchaseAgreementGateway: notePurchaseAgreementGateway,
		userService:                  userService,
	}
}
