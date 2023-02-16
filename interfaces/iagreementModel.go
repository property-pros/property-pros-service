package interfaces

import (
	"github.com/vireocloud/property-pros-service/interop"
)

type IAgreementModel interface {
	IBaseModel[interop.NotePurchaseAgreement]
	Save() (IAgreementModel, error)
	LoadDocument() (IAgreementModel, error)
}
