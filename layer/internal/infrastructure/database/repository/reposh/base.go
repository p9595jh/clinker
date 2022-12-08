package reposh

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/p9595jh/fpgo"
	"github.com/p9595jh/transform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// to represent none-dao entity
type Nil struct{}

// D: DAO
// E: Entity
type BaseRepository[D, E any] struct {
	daoExists   bool
	collection  *mongo.Collection
	transformer transform.Transformer
}

func NewBaseRepository[D, E any](collection *mongo.Collection, transformer transform.Transformer) *BaseRepository[D, E] {
	return &BaseRepository[D, E]{
		daoExists:   !(reflect.TypeOf(new(D)) == reflect.TypeOf(&Nil{})),
		collection:  collection,
		transformer: transformer,
	}
}

func (r *BaseRepository[D, E]) Insert(e *E, opts ...*options.InsertOneOptions) (any, error) {
	var doc any
	if r.daoExists {
		doc = new(D)
		if err := r.transformer.Mapping(e, doc); err != nil {
			return nil, err
		}
	} else {
		doc = e
	}
	if res, err := r.collection.InsertOne(context.Background(), doc, opts...); err != nil {
		return nil, err
	} else {
		return res.InsertedID, nil
	}
}

func (r *BaseRepository[D, E]) Find(filter any, opts ...*options.FindOptions) ([]E, error) {
	if filter == nil {
		filter = bson.M{}
	}
	if res, err := r.collection.Find(context.Background(), filter, opts...); err != nil {
		return nil, err
	} else {
		if r.daoExists {
			var daos []*D
			if err := res.All(context.Background(), &daos); err != nil {
				return nil, err
			} else {
				errs := []string{}
				es := fpgo.Pipe[[]*D, []E](
					daos,
					fpgo.Map(func(i int, d **D) *E {
						e := new(E)
						if err := r.transformer.Mapping(*d, e); err != nil {
							errs = append(errs, err.Error())
						}
						return e
					}),
				)
				if len(errs) > 0 {
					return nil, errors.New(strings.Join(errs, "; "))
				} else {
					return es, nil
				}
			}
		} else {
			var data []E
			err := res.All(context.Background(), &data)
			return data, err
		}
	}
}

func (r *BaseRepository[D, E]) FindOne(filter any, opts ...*options.FindOneOptions) (*E, error) {
	if filter == nil {
		filter = bson.M{}
	}
	res := r.collection.FindOne(context.Background(), filter, opts...)
	if r.daoExists {
		var (
			d = new(D)
			e = new(E)
		)
		err := res.Decode(d)
		if err != nil {
			return nil, err
		}
		err = r.transformer.Mapping(d, e)
		if err != nil {
			return nil, err
		}
		return e, nil
	} else {
		e := new(E)
		if err := res.Decode(e); err != nil {
			return nil, err
		}
		return e, nil
	}
}

func (r *BaseRepository[D, E]) Update(filter, updates any, opts ...*options.FindOneAndUpdateOptions) error {
	if filter == nil {
		filter = bson.M{}
	}
	return r.collection.FindOneAndUpdate(
		context.Background(),
		filter,
		updates,
		opts...,
	).Err()
}

func (r *BaseRepository[D, E]) Delete(filter any, opts ...*options.FindOneAndDeleteOptions) error {
	if filter == nil {
		filter = bson.M{}
	}
	return r.collection.FindOneAndDelete(
		context.Background(),
		filter,
		opts...,
	).Err()
}
