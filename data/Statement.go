package data

import "time"

type Statement struct {
	CreatedOn time.Time `json:"created_on" gorm:"column:created_on;autoCreateTime"`
	UpdatedOn time.Time `json:"updated_on" gorm:"column:updated_on;autoUpdateTime"`

	Id              string `json:"id,omitempty" gorm:"primaryKey"`
	UserId          string `json:"userid" gorm:"column:user_id"`
	EmailAddress    string `json:"email_address" gorm:"column:email_address"`
	StartPeriodDate time.Time `json:"start_period_date" gorm:"column:start_period_date"`
	EndPeriodDate   time.Time `json:"end_period_date" gorm:"column:end_period_date"`
	Balance         string `json:"balance" gorm:"column:balance"`
	TotalIncome     string `json:"total_income" gorm:"column:total_income"`
	Principle       string `json:"principle" gorm:"column:principle"`
}

func (n *Statement) GetId() string {
	return n.Id
}

func (u *Statement) TableName() string {
	return "statements"
}
