package data

type User struct {
	Id             string `json:"id,omitempty" gorm:"primaryKey"`
	FirstName      string `json:"first_name,omitempty" gorm:"column:first_name"`
	LastName       string `json:"last_name,omitempty" gorm:"column:last_name"`
	DateOfBirth    string `json:"date_of_birth,omitempty" gorm:"column:date_of_birth"`
	EmailAddress   string `json:"email_address,omitempty" gorm:"column:email_address"`
	Password       string `json:"password,omitempty" gorm:"column:password"`
	HomeAddress    string `json:"home_address,omitempty" gorm:"column:home_address"`
	PhoneNumber    string `json:"phone_number,omitempty" gorm:"column:phone_number"`
	SocialSecurity string `json:"social_security,omitempty" gorm:"column:social_security"`
	CreatedOn      string `json:"created_on,omitempty" gorm:"column:created_on"`
}

func (u *User) GetId() string {
	return u.Id
}

func (u *User) TableName() string {
	return "users"
}
