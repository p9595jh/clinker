package entity

import (
	"clinker-backend/internal/infrastructure/database/entity/gormem"
	"database/sql"
)

type User struct {
	gormem.BaseEntity `gorm:"embedded"`
	Id                string       `gorm:"type:char(30);index:idx_user_id;not null"`
	Password          string       `gorm:"type:char(60);not null" transform:"map:-"`
	Nickname          string       `gorm:"type:char(20);index:idx_user_nickname;not null"`
	Address           string       `gorm:"type:char(40);charset latin1;default:0000000000000000000000000000000000000000;not null" transform:"add0x"`
	Authority         uint8        `gorm:"type:tinyint;default:0;not null"` // 0: user, 1: administrator
	StopReason        string       `gorm:"type:char(30);null"`
	StopUntil         sql.NullTime `gorm:"type:datetime(6) null" transform:"timeFormat"`

	Vestiges   []Vestige   `transform:"vestigesE2R"`
	Appraisals []Appraisal `transform:"appraisalsE2R"`
}
