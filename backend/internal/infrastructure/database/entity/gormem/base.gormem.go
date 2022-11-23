package gormem

import (
	"database/sql"
	"time"
)

type BaseEntity struct {
	CreatedAt time.Time `gorm:"type:datetime(6) default CURRENT_TIMESTAMP(6);not null" transform:"timeFormat,map:CreatedAt"`
	UpdatedAt time.Time `gorm:"type:datetime(6) default CURRENT_TIMESTAMP(6) not null on update CURRENT_TIMESTAMP(6)" transform:"timeFormat,map:-"`
}

type BaseEntityAutoId struct {
	BaseEntity `gorm:"embedded"`
	Id         uint `gorm:"primarykey;type:int(11)" transform:"map:Id"`
}

type BaseEntityUuid struct {
	BaseEntity `gorm:"embedded"`
	Id         string `gorm:"type:char(36);primaryKey" transform:"map:Id"`
}

type BaseEntityTxHash struct {
	BaseEntity `gorm:"embedded"`
	TxHash     string `gorm:"type:char(64);primaryKey" transform:"map:TxHash"`
}

type DeletedAt struct {
	DeletedAt sql.NullTime `gorm:"type:datetime(6) null" transform:"map:-"`
}
