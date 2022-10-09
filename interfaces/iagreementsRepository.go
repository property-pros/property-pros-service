package interfaces

import "github.com/vireocloud/property-pros-service/interop"

type IAgreementsRepository interface {
	IRepository[interop.NotePurchaseAgreement]
}
