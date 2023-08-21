package agreements

import (
	"context"
	"fmt"
	"log"
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

func (g *NotePurchaseAgreementGateway) Getall(ctx context.Context, usrID string) ([]*interop.NotePurchaseAgreement, error) {
	agreements := make([]*interop.NotePurchaseAgreement, 0)

	npas := g.npaRepository.Query(&data.NotePurchaseAgreement{
		UserId: usrID,
	})
	log.Printf("id: %+v", usrID)
	log.Printf("npas: %+v", npas)
	for _, npa := range npas {
		agreements = append(agreements, &interop.NotePurchaseAgreement{
			Id:        npa.Id,
			CreatedOn: npa.CreatedOn.Format(time.RFC3339),
		})
	}

	return agreements, nil
}

func (g *NotePurchaseAgreementGateway) FindOne(ctx context.Context, payload interfaces.IModelPayload) (npaRecord *interop.NotePurchaseAgreement, docURL string, err error) {
	npaReqModel := data.NotePurchaseAgreement{
		Id: payload.GetId(),
	}

	npa, err := g.npaRepository.FindOne(&npaReqModel)
	if err != nil {
		return nil, "", fmt.Errorf("error finding note purchase agreement: %w", err)
	}

	usr, err := g.userRepository.FindOne(&data.User{
		Id: npa.UserId,
	})
	if err != nil {
		return nil, "", fmt.Errorf("error finding user: %w", err)
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
		User: &interop.User{
			Id:           usr.Id,
			EmailAddress: usr.EmailAddress,
			Password:     usr.Password,
		},
	}, npa.DocURL, nil

}
