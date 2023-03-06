package agreements

import (
	"context"
	"fmt"
	"log"

	"github.com/vireocloud/property-pros-service/documents"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/users"
)

type NotePurchaseAgreementService struct {
	factory                      interfaces.INotePurchaseAgreementModelFactory
	notePurchaseAgreementGateway *NotePurchaseAgreementGateway
	usersGateway                 *users.UsersGateway
	documentContentService       *documents.DocumentContentService
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
	return service.notePurchaseAgreementGateway.FindOne(ctx, payload)
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreements(ctx context.Context) ([]interfaces.IAgreementModel, error) {
	return service.notePurchaseAgreementGateway.Getall(ctx)
}

func (service *NotePurchaseAgreementService) Save(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {

	user, err := service.usersGateway.SaveUser(ctx, agreement)

	if err != nil {
		return nil, err
	}

	agreement.User.Id = user.Id

	docURL, err := service.documentContentService.CreateAndSaveNotePurchaseAgreementDoc(ctx, agreement)
	if err != nil {
		return nil, fmt.Errorf("failed to generate doc content, err: %w", err)
	}

	resultAgreement, err := service.notePurchaseAgreementGateway.SaveNotePurchaseAgreement(ctx, agreement, docURL)

	if err != nil {
		return nil, err
	}

	resultAgreement.User.Id = user.Id

	return resultAgreement, nil
}

func NewNotePurchaseAgreementService(
	factory interfaces.INotePurchaseAgreementModelFactory,
	npag *NotePurchaseAgreementGateway,
	usersGateway *users.UsersGateway,
	documentContentService *documents.DocumentContentService,
) interfaces.IAgreementsService {
	log.Printf("factory: %+#v \n\n", factory)
	return &NotePurchaseAgreementService{
		factory:                      factory,
		notePurchaseAgreementGateway: npag,
		usersGateway:                 usersGateway,
		documentContentService:       documentContentService,
	}
}

var _ interfaces.IAgreementsService = (*NotePurchaseAgreementService)(nil)
