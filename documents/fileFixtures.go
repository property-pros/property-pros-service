package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

func InitFileFixtures(notePurchaseAgreementsService interfaces.IAgreementsService) error {

	notePurchaseAgreement, err := notePurchaseAgreementsService.GetNotePurchaseAgreement(context.TODO(), &interop.NotePurchaseAgreement{Id: "test-npa-1"})

	if err != nil {
		return err
	}
	
	_, err = notePurchaseAgreementsService.Save(context.TODO(), notePurchaseAgreement)

	if err != nil {
		return err
	}

	return nil
}