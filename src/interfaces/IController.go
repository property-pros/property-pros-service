package interfaces

type IController interface {
	IBaseModel
	Save() IGroup
	Delete() IGroup
}
