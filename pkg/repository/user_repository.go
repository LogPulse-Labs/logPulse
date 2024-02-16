package repository

import (
	"context"
	"log-pulse/models"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	Collection *mongo.Collection
}

type UserRepository interface {
	CreateOne(user *models.User, ctx ...*mongo.SessionContext) (*models.User, error)
	All() (*[]models.User, error)
	FindOne(ID interface{}) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(ID string) error
	GetCollection() *mongo.Collection
}

func (r *userRepository) CreateOne(user *models.User, ctx ...*mongo.SessionContext) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = ""
	user.EmailVerified = false

	hashPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashPassword

	var contextValue context.Context

	if len(ctx) > 0 && ctx[0] != nil {
		contextValue = *ctx[0]
	} else {
		contextValue = context.Background()
	}

	result, err := r.Collection.InsertOne(contextValue, user)

	if err != nil {
		return nil, err
	}

	createdUser, findError := r.FindOne(result.InsertedID)
	if findError != nil {
		return nil, findError
	}

	return createdUser, nil
}

// Delete implements UserRepository.
func (r *userRepository) Delete(ID string) error {
	panic("unimplemented")
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var result models.User

	if err := r.Collection.FindOne(
		context.Background(),
		bson.D{{Key: "email", Value: email}},
		options.FindOne(),
	).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *userRepository) FindOne(ID interface{}) (*models.User, error) {
	var result models.User

	if err := r.Collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: ID}}).
		Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get implements UserRepository.
func (r *userRepository) All() (*[]models.User, error) {
	panic("unimplemented")
}

// Update implements UserRepository.
func (r *userRepository) Update(model *models.User) (*models.User, error) {
	panic("unimplemented")
}

func (r *userRepository) GetCollection() *mongo.Collection {
	return r.Collection
}

func NewUserRepository() UserRepository {
	return &userRepository{
		Collection: database.DB.Collection("users"),
	}
}
