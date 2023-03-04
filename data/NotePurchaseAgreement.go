package data

import "time"

type NotePurchaseAgreement struct {
	Id             string    `json:"id,omitempty" gorm:"primaryKey"`
	FundsCommitted uint64    `json:"funds_committed,omitempty" gorm:"column:funds_commited"`
	UserId         string    `json:"userid" bson:"userid" gorm:"column:user_id"`
	CreatedOn      time.Time `json:"created_on" gorm:"column:created_on;autoCreateTime"`
	UpdatedOn      time.Time `json:"updated_on" gorm:"column:updated_on;autoUpdateTime"`
}

func (n *NotePurchaseAgreement) GetId() string {
	return n.Id
}

func (u *NotePurchaseAgreement) TableName() string {
	return "note_purchase_agreements"
}
