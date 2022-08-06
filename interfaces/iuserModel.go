package interfaces

type IUserModel interface {
	GetId() uint32
	GetFirstName() string
	GetLastName() string
	GetDateOfBirth() string
	GetHomeAddress() string
	GetEmailAddress() string
	GetPhoneNumber() string
	GetSocialSecurity() string
}
