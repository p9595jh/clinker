package dbconn

import (
	"clinker-backend/common/config"
	"clinker-backend/internal/infrastructure/database/entity"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Schema   string
}

func Connect(dbconf *DatabaseConfig) (*gorm.DB, error) {
	dialector := mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=latin1&parseTime=True",
		dbconf.User,
		dbconf.Password,
		dbconf.Host,
		dbconf.Port,
		dbconf.Schema,
	))

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &customLogger{},
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		panic(err)
	}

	if config.V.GetBool("db.sync") {
		if err = db.AutoMigrate(
			&entity.User{},
			&entity.Vestige{},
			&entity.Appraisal{},
		); err != nil {
			Close(db)
			return nil, err
		}
	}

	return db, nil
}

func Close(db *gorm.DB) error {
	d, err := db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}
