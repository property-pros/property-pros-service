package agreements

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/constants"
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

	agreementModelData := data.NotePurchaseAgreement{
		Id:             uuid.New().String(),
		FundsCommitted: agreement.FundsCommitted,
		UserId:         agreement.User.Id,
	}

	agreementSaved, err := g.npaRepository.Save(&agreementModelData)
	if err != nil {
		return nil, err
	}

	agreement.Id = agreementSaved.Id
	agreement.CreatedOn = agreementSaved.CreatedOn.Format(time.RFC3339)

	return agreement, nil
}

func (g *NotePurchaseAgreementGateway) SaveNotePurchaseAgreement(ctx context.Context, agreement *interop.NotePurchaseAgreement, docURL string) (*interop.NotePurchaseAgreement, error) {
	agreementModelData := data.NotePurchaseAgreement{
		Id:             uuid.New().String(),
		FundsCommitted: agreement.FundsCommitted,
		UserId:         agreement.User.Id,
		DocURL:         docURL,
	}

	agreementSaved, err := g.npaRepository.Save(&agreementModelData)
	if err != nil {
		return nil, err
	}

	agreement.Id = agreementSaved.Id
	agreement.CreatedOn = agreementSaved.CreatedOn.Format(time.RFC3339)

	return agreement, nil
}

func (g *NotePurchaseAgreementGateway) Getall(ctx context.Context) ([]*interop.NotePurchaseAgreement, error) {
	agreements := make([]*interop.NotePurchaseAgreement, 0)
	fmt.Println("ctx.Value(constants.UserIdKey)")
	fmt.Println(ctx.Value(constants.UserIdKey))
	usrID := fmt.Sprintf("%v", ctx.Value(constants.UserIdKey))
	if usrID == "" {
		return nil, errors.New("unresolved userid")
	}

	npas := g.npaRepository.Query(&data.NotePurchaseAgreement{
		UserId: usrID,
	})
	fmt.Println(npas)

	for _, npa := range npas {
		agreements = append(agreements, &interop.NotePurchaseAgreement{
			Id:        npa.Id,
			CreatedOn: npa.CreatedOn.Format(time.RFC3339),
		})
	}
	fmt.Println(agreements)

	return agreements, nil
}

func (g *NotePurchaseAgreementGateway) FindOne(ctx context.Context, payload interfaces.IModelPayload) (npaRecord *interop.NotePurchaseAgreement, docURL string, err error) {
	npaReqModel := data.NotePurchaseAgreement{
		Id: payload.GetId(),
	}

	npa, err := g.npaRepository.FindOne(&npaReqModel)
	if err != nil {
		return nil, "", err
	}

	usr, err := g.userRepository.FindOne(&data.User{
		Id: npa.UserId,
	})
	if err != nil {
		return nil, "", err
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
	}, npa.DocURL, nil
}
