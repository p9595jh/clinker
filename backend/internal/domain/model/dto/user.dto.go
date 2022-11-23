package dto

import "time"

type UserPanginationDto struct {
	Skip int `json:"skip" mapstructure:"skip" validate:"required,numeric,gte=0"`
	Take int `json:"take" mapstructure:"take" validate:"required,numeric,gte=0"`
}

type UserIdDto struct {
	UserId string `json:"userId" mapstructure:"userId" validate:"required"`
}

type UserStopDtom struct {
	Reason string `json:"reason" mapstructure:"reason" validate:"required"`
	Date   string `json:"date" mapstructure:"date" validate:"required,date" transform:"timeFormat"`
}

type UserStopDto struct {
	Reason string
	Date   time.Time
}
