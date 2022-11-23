package dto

type AppraisalTxHashDto struct {
	TxHash string `json:"txHash" mapstructure:"txHash" validate:"required,txHash" transform:"hex"`
}

type AppraisalUserIdDto struct {
	UserId string `json:"userId" mapstructure:"userId" validate:"required"`
}

type AppraisalPanginationDto struct {
	Skip int `json:"skip" mapstructure:"skip" validate:"required,numeric,gte=0"`
	Take int `json:"take" mapstructure:"take" validate:"required,numeric,gte=0"`
}
