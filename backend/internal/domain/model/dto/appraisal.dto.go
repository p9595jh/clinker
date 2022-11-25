package dto

type AppraisalDto struct {
	Value  int64  `json:"value" validate:"required,numeric,min=-50,max=50" example:"30"`
	NextId string `json:"vestige" validate:"required,txHash" transform:"hex"`
}
