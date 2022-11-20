package entity

import "clinker-backend/internal/infrastructure/database/entity/gormem"

type Vestige struct {
	gormem.BaseEntityTxHash `gorm:"embedded"`
	Parent                  string `gorm:"type:char(64);index:idx_vestige_parent;charset latin1;default:0000000000000000000000000000000000000000000000000000000000000000;not null"`
	Head                    string `gorm:"type:char(64);index:idx_vestige_head;charset latin1;default:0000000000000000000000000000000000000000000000000000000000000000;not null"`
	Next                    string `gorm:"type:char(64);index:idx_vestige_next;charset latin1;default:0000000000000000000000000000000000000000000000000000000000000000;not null"`
	Title                   string `gorm:"type:text;not null"`
	Content                 string `gorm:"type:text;not null"`
	Hit                     int64  `gorm:"type:int;default:0;not null"`
	Confirmed               bool   `gorm:"type:bool;default:false;not null"`

	UserId     string
	User       User
	Appraisals []Appraisal

	Children []Vestige `gorm:"-"`
}
