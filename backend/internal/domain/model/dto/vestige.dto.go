package dto

type VestigePanginationDto struct {
	Skip int `json:"skip" mapstructure:"skip" validate:"required,numeric,gte=0"`
	Take int `json:"take" mapstructure:"take" validate:"required,numeric,gte=0"`
}

type VestigeTxHashDto struct {
	TxHash string `json:"txHash" mapstructure:"txHash" validate:"required,txHash" transform:"hex"`
}
