package entity

import "clinker-backend/internal/infrastructure/database/entity/gormem"

type Appraisal struct {
	gormem.BaseEntityTxHash `gorm:"embedded"`
	Value                   int64 `gorm:"type:int;not null"`
	Confirmed               bool  `gorm:"type:bool;default:false;not null"`

	VestigeId string // head
	Vestige   Vestige
	NextId    string // actual target
	Next      Vestige
	UserId    string
	User      User
}
