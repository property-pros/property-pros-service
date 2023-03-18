package agreements

import (
	"context"
	"fmt"
	"log"

	"github.com/vireocloud/property-pros-service/documents"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementService struct {
	factory                      interfaces.INotePurchaseAgreementModelFactory
	notePurchaseAgreementGateway *NotePurchaseAgreementGateway
	usersGateway                 interfaces.IUsersGateway
	documentContentService       *documents.DocumentContentService
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreementDocContent(ctx context.Context, payload interfaces.IModelPayload) ([]byte, error) {

	_, docURL, err := service.notePurchaseAgreementGateway.FindOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	content, err := service.documentContentService.GetDocumentContent(ctx, docURL)
	fmt.Printf("content is %v\n", content)
	fmt.Printf("content is %v\n", string(content))
	return content, err
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreement(ctx context.Context, payload interfaces.IModelPayload) (*interop.NotePurchaseAgreement, error) {
	npa, _, err := service.notePurchaseAgreementGateway.FindOne(ctx, payload)
	return npa, err
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
	usersGateway interfaces.IUsersGateway,
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
