package agreements

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/users"
)

type NotePurchaseAgreementService struct {
	factory                      interfaces.INotePurchaseAgreementModelFactory
	notePurchaseAgreementGateway *NotePurchaseAgreementGateway
	userGateway                  *users.UsersGateway
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreementDocContent(ctx context.Context, payload interfaces.IModelPayload) ([]byte, error) {
	model, err := service.GetNotePurchaseAgreementModel(ctx, payload)

	if err != nil {
		return nil, err
	}

	model, err = model.LoadDocument()

	if err != nil {
		return nil, err
	}

	return model.GetPayload().GetFileContent(), nil
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreementModel(ctx context.Context, payload interfaces.IModelPayload) (interfaces.IAgreementModel, error) {

	// model, err := service.factory.NewPurchaseAgreementModel(ctx, payload.(*interop.NotePurchaseAgreement))

	// if err != nil {
	// 	return nil, err
	// }

	// model, err = service.notePurchaseAgreementGateway.FindOne(ctx, model)

	// if err != nil {
	// 	return nil, err
	// }

	// return model, nil
	return &NotePurchaseAgreementModel{}, nil
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreement(ctx context.Context, payload interfaces.IModelPayload) (*interop.NotePurchaseAgreement, error) {

	model, err := service.GetNotePurchaseAgreementModel(ctx, payload)

	if err != nil {
		return nil, err
	}

	return model.GetPayload(), nil
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreements(ctx context.Context) ([]interfaces.IAgreementModel, error) {
	return service.notePurchaseAgreementGateway.Getall(ctx)
}

func (service *NotePurchaseAgreementService) Save(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {
	userData := data.User{
		Id:             uuid.New().String(),
		FirstName:      agreement.GetFirstName(),
		LastName:       agreement.GetLastName(),
		DateOfBirth:    agreement.GetDateOfBirth(),
		EmailAddress:   agreement.User.GetEmailAddress(),
		Password:       agreement.User.GetPassword(),
		HomeAddress:    agreement.GetHomeAddress(),
		PhoneNumber:    agreement.GetPhoneNumber(),
		SocialSecurity: agreement.GetSocialSecurity(),
		CreatedOn:      agreement.GetCreatedOn(),
	}

	user, err := service.userGateway.SaveUser(userData)
	if err != nil {
		return nil, err
	}

	agreementModelData := data.NotePurchaseAgreement{
		Id:             uuid.New().String(),
		FundsCommitted: agreement.FundsCommitted,
		UserId:         user.Id,
	}

	_, err = service.notePurchaseAgreementGateway.SaveNotePurchaseAgreement(ctx, agreementModelData)
	if err != nil {
		return nil, err
	}

	return agreement, nil
}

func NewNotePurchaseAgreementService(
	factory interfaces.INotePurchaseAgreementModelFactory,
	npag *NotePurchaseAgreementGateway,
	ug *users.UsersGateway) interfaces.IAgreementsService {
	log.Printf("factory: %+#v \n\n", factory)
	return &NotePurchaseAgreementService{
		factory:                      factory,
		notePurchaseAgreementGateway: npag,
		userGateway:                  ug,
	}
}

var _ interfaces.IAgreementsService = (*NotePurchaseAgreementService)(nil)
