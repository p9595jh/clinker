package dto

type QueryPaginationDto struct {
	Page *int `json:"page" mapstructure:"page" query:"page" validate:"required,numeric,gte=0" example:"0"`
	Take int  `json:"take" mapstructure:"take" query:"take" validate:"required,numeric,gt=0" example:"10"`
}

type ParamTxHashDto struct {
	TxHash string `json:"txHash" mapstructure:"txHash" validate:"required,txHash" transform:"hex" example:"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"`
}

type ParamAddressDto struct {
	Address string `json:"address" mapstructure:"address" validate:"required,ethAddr" transform:"hex" example:"0x1234567890abcdef1234567890abcdef12345678"`
}

type ParamUserIdDto struct {
	UserId string `json:"userId" mapstructure:"userId" validate:"required" example:"user123"`
}
