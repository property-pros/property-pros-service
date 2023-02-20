package data

type NotePurchaseAgreement struct {
	Id             string `json:"id,omitempty" gorm:"primaryKey"`
	FundsCommitted uint64 `json:"funds_committed,omitempty" gorm:"column:funds_commited"`
	UserId         string `json:"userid" bson:"userid" gorm:"column:user_id"`
}

func (n *NotePurchaseAgreement) GetId() string {
	return n.Id
}

func (u *NotePurchaseAgreement) TableName() string {
	return "note_purchase_agreements"
}
