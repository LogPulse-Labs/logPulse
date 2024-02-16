package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(model *interface{}) (*interface{}, error)
	Get() (*[]any, error)
	FindOne(filter interface{}, options *options.FindOneOptions) (any, error)
	Update(model *any) (*any, error)
	Delete(ID string) error
}

type repository struct {
	Collection *mongo.Collection
}

// Create implements Repository.
func (r *repository) Create(model *interface{}) (*interface{}, error) {
	result, err := r.Collection.InsertOne(context.Background(), model)

	if err != nil {
		return nil, err
	}

	created, findError := r.FindOne(result.InsertedID, nil)
	if findError != nil {
		return nil, findError
	}

	return &created, nil
}

// Delete implements Repository.
func (r *repository) Delete(ID string) error {
	panic("unimplemented")
}

// FindOne implements Repository.
func (r *repository) FindOne(filter interface{}, options *options.FindOneOptions) (any, error) {
	var result any

	err := r.Collection.FindOne(context.Background(), filter, options).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Get implements Repository.
func (r *repository) Get() (*[]any, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (*repository) Update(model *any) (*any, error) {
	panic("unimplemented")
}

func NewBaseRepository(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
