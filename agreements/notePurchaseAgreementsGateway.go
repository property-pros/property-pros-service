package agreements

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementGateway struct {
	npaRepository  interfaces.IRepository[data.NotePurchaseAgreement]
	userRepository interfaces.IRepository[data.User]
	factory        interfaces.INotePurchaseAgreementModelFactory
}

func NewNotePurchaseAgreementGateway(
	npaRepository interfaces.IRepository[data.NotePurchaseAgreement],
	userRepository interfaces.IRepository[data.User],
	factory interfaces.INotePurchaseAgreementModelFactory) *NotePurchaseAgreementGateway {
	return &NotePurchaseAgreementGateway{
		npaRepository:  npaRepository,
		userRepository: userRepository,
		factory:        factory,
	}
}

func (g *NotePurchaseAgreementGateway) SaveUserAndNotePurchaseAgreement(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {
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
	}

	user, err := g.userRepository.Save(&userData)
	if err != nil {
		return nil, err
	}

	agreementModelData := data.NotePurchaseAgreement{
		Id:             uuid.New().String(),
		FundsCommitted: agreement.FundsCommitted,
		UserId:         user.Id,
	}

	agreementSaved, err := g.npaRepository.Save(&agreementModelData)
	if err != nil {
		return nil, err
	}

	agreement.Id = agreementSaved.Id
	agreement.CreatedOn = agreementSaved.CreatedOn.Format(time.RFC3339)

	return agreement, nil
}

func (g *NotePurchaseAgreementGateway) SaveNotePurchaseAgreement(ctx context.Context, agreement data.NotePurchaseAgreement) (data.NotePurchaseAgreement, error) {
	agreement.Id = uuid.New().String()

	_, err := g.npaRepository.Save(&agreement)
	if err != nil {
		return data.NotePurchaseAgreement{}, err
	}

	return agreement, nil
}

func (g *NotePurchaseAgreementGateway) Getall(ctx context.Context) ([]interfaces.IAgreementModel, error) {
	return nil, nil
	// results := g.repository.Query(nil)
	// models := []interfaces.IAgreementModel{}

	// for _, result := range results {
	// 	model, err := g.factory.NewPurchaseAgreementModel(ctx, result)

	// 	if err != nil {
	// 		return models, err
	// 	}

	// 	models = append(models, model)
	// }

	// return models, nil
}

func (g *NotePurchaseAgreementGateway) FindOne(ctx context.Context, payload interfaces.IModelPayload) (*interop.NotePurchaseAgreement, error) {
	npaReqModel := data.NotePurchaseAgreement{
		Id: payload.GetId(),
	}

	npa, err := g.npaRepository.FindOne(&npaReqModel)
	if err != nil {
		return nil, err
	}

	usr, err := g.userRepository.FindOne(&data.User{
		Id: npa.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &interop.NotePurchaseAgreement{
		Id:             npa.Id,
		FirstName:      usr.FirstName,
		LastName:       usr.LastName,
		DateOfBirth:    usr.DateOfBirth,
		HomeAddress:    usr.HomeAddress,
		PhoneNumber:    usr.PhoneNumber,
		SocialSecurity: usr.SocialSecurity,
		FundsCommitted: npa.FundsCommitted,
		CreatedOn:      npa.CreatedOn.Format(time.RFC3339),
	}, nil
}

func (g *NotePurchaseAgreementGateway) FindOne1(ctx context.Context, model data.NotePurchaseAgreement) (*data.NotePurchaseAgreement, error) {
	return g.npaRepository.FindOne(&model)
}
