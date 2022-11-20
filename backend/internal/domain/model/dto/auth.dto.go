package dto

type AuthLoginDto struct {
	Id       string `json:"id" mapstructure:"id" validate:"required,min=5,max=30" example:"user123"`
	Password string `json:"password" mapstructure:"password" validate:"required,min=8,max=30" example:"123123123"`
}

type AuthRegisterDtom struct {
	Id       string `json:"id" validate:"required,min=5,max=30" example:"user123"`
	Password string `json:"password" validate:"required,min=8,max=30" example:"123123123"`
	NickName string `json:"name" validate:"required,min=2,max=20" example:"john"`
	Address  string `json:"address" validate:"required,ethAddr" exmaple:"1234567890abcdef1234567890abcdef12345678" transform:"hex"`
}

type AuthRegisterDto struct {
	Id       string
	Password string
	NickName string
	Address  string
}
