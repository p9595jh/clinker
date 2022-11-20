package reposh

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[E any] interface {
	Find(findOption *FindOption[E]) (*[]E, error)
	FindOne(findOption *FindOption[E]) (*E, error)
	Save(e *E) (*E, error)
	Update(where *EntityParts[E], update *EntityParts[E]) *gorm.DB
	UpdateAndCall(where *EntityParts[E], update *EntityParts[E], callbacks ...func(result *gorm.DB)) error
	Delete(where *EntityParts[E], id interface{}) error
	DeleteById(id any) error
	Model() *gorm.DB
	Exec(tx *gorm.DB) (*[]E, error)
	ExecOne(tx *gorm.DB) (*E, error)
}

type baseRepository[E any] struct {
	model *E
	DB    *gorm.DB
}

func NewRepository[E any](db *gorm.DB) BaseRepository[E] {
	return &baseRepository[E]{
		model: new(E),
		DB:    db,
	}
}

func (r *baseRepository[E]) Model() *gorm.DB {
	return r.DB.Model(r.model)
}

func (r *baseRepository[E]) where(tx *gorm.DB, findOption *FindOption[E]) *gorm.DB {
	if findOption == nil {
		return tx
	}

	switch {
	case findOption.Where.Entities != nil:
		conditions := *findOption.Where.Entities
		tx = tx.Where(conditions[0])
		for i := 1; i < len(conditions); i++ {
			tx = tx.Or(conditions[i])
		}
	case findOption.Where.Entity != nil:
		tx = tx.Where(findOption.Where.Entity)
	case findOption.Where.EntitiesFn != nil:
		args := make([]E, 2)
		findOption.Where.EntitiesFn(&args)
		tx = tx.Where(args[0])
		for i := 1; i < len(args); i++ {
			tx = tx.Or(args[i])
		}
	case findOption.Where.EntityFn != nil:
		arg := new(E)
		findOption.Where.EntityFn(arg)
		tx = tx.Where(arg)
	case len(findOption.Where.Raw) > 0:
		tx = tx.Where(findOption.Where.Raw[0], findOption.Where.Raw[1:]...)
	}

	if len(findOption.Preload) > 0 {
		for _, p := range findOption.Preload {
			tx = tx.Preload(p)
		}
	}
	if findOption.Limit > 0 {
		tx = tx.Limit(findOption.Limit)
	}
	if findOption.Offset > 0 {
		tx = tx.Offset(findOption.Offset)
	}
	return tx
}

func (r *baseRepository[E]) build(findOption *FindOption[E]) *gorm.DB {
	tx := r.Model()
	if findOption == nil {
		return tx
	}

	if len(findOption.Select) > 0 {
		tx = tx.Select(strings.Join(findOption.Select, ","))
	}

	tx = r.where(tx, findOption)
	if findOption.Order.Column != "" {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: findOption.Order.Column},
			Desc:   findOption.Order.Desc,
		})
	}
	return tx
}

func (r *baseRepository[E]) Find(findOption *FindOption[E]) (*[]E, error) {
	instances := new([]E)
	tx := r.build(findOption).Find(instances)
	return instances, tx.Error
}

func (r *baseRepository[E]) FindOne(findOption *FindOption[E]) (*E, error) {
	instance := new(E)
	tx := r.build(findOption).Take(instance)
	return FilteredRecord(instance, tx.Error)
}

func (r *baseRepository[E]) Save(e *E) (*E, error) {
	err := r.DB.Save(e).Error
	return e, err
}

func (r *baseRepository[E]) Update(where *EntityParts[E], update *EntityParts[E]) *gorm.DB {
	tx := r.Model()
	tx = r.where(tx, &FindOption[E]{Where: *where})
	switch {
	case update.EntityFn != nil:
		arg := new(E)
		update.EntityFn(arg)
		return tx.Updates(arg)
	case update.EntitiesFn != nil:
		args := make([]E, 2)
		update.EntitiesFn(&args)
		return tx.Updates(&args[0])
	}
	return nil
}

func (r *baseRepository[E]) UpdateAndCall(where *EntityParts[E], update *EntityParts[E], callbacks ...func(result *gorm.DB)) error {
	tx := r.Update(where, update)
	if tx.Error != nil {
		return tx.Error
	}
	for _, callback := range callbacks {
		callback(tx)
	}
	return nil
}

func (r *baseRepository[E]) Delete(where *EntityParts[E], id interface{}) error {
	return r.where(r.Model(), &FindOption[E]{Where: *where}).Delete(id).Error
}

func (r *baseRepository[E]) DeleteById(id any) error {
	return r.DB.Delete(id).Error
}

func (r *baseRepository[E]) Exec(tx *gorm.DB) (*[]E, error) {
	var items []E
	err := tx.Find(&items).Error
	return &items, err
}

func (r *baseRepository[E]) ExecOne(tx *gorm.DB) (*E, error) {
	var item E
	err := tx.Take(&item).Error
	return &item, err
}
