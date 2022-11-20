package entity

import "clinker-backend/internal/infrastructure/database/entity/gormem"

type User struct {
	gormem.BaseEntity `gorm:"embedded"`
	Id                string `gorm:"type:char(30);index:idx_user_id;not null"`
	Password          string `gorm:"type:char(60);not null"`
	Nickname          string `gorm:"type:char(20);index:idx_user_nickname;not null"`
	Address           string `gorm:"type:char(40);charset latin1;default:0000000000000000000000000000000000000000;not null"`

	Vestiges   []Vestige
	Appraisals []Appraisal
}
