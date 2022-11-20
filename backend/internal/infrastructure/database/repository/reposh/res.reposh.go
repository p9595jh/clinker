package reposh

import (
	"errors"

	"gorm.io/gorm"
)

func FilteredRecord[E any](record *E, err error) (*E, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return record, nil
}
