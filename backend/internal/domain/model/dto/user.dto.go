package dto

import "time"

type UserStopDtom struct {
	Reason string `json:"reason" mapstructure:"reason" validate:"required" example:"no specific reason"`
	Date   string `json:"date" mapstructure:"date" validate:"required,date" transform:"timeFormat" example:"2100-01-01"`
}

type UserStopDto struct {
	Reason string
	Date   time.Time
}

type UserDto struct {
	Id       string `json:"id" validate:"required,max=30" example:"user123"`
	Password string `json:"password" validate:"required" example:"123123123"`
	Nickname string `json:"nickname" validate:"required,max=20" example:"nick"`
	Address  string `json:"address" validate:"required,ethAddr" example:"0x1234567890abcdef1234567890abcdef12345678" transform:"hex"`
}
