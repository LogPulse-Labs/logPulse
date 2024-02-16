package repository

import (
	"context"
	"log-pulse/models"
	"log-pulse/platform/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrganizationRepository interface {
	CreateOne(organization *models.Organization) (*models.Organization, error)
	All() (*[]models.Organization, error)
	FindOne(filter interface{}) (*models.Organization, error)
	FindByUser(userId string) (*models.Organization, error)
	Update(organization *models.Organization) (*models.Organization, error)
	Delete(ID string) error
	GetCollection() *mongo.Collection
}

type organizationRepository struct {
	Collection *mongo.Collection
}

// Create implements Repository.
func (r *organizationRepository) CreateOne(organization *models.Organization) (*models.Organization, error) {
	organization.CreatedAt = time.Now()
	organization.UpdatedAt = time.Now()
	organization.ID = ""

	result, err := r.Collection.InsertOne(context.Background(), organization)

	if err != nil {
		return nil, err
	}

	created, findError := r.FindOne(result.InsertedID)
	if findError != nil {
		return nil, findError
	}

	return created, nil
}

func (r *organizationRepository) Delete(ID string) error {
	panic("unimplemented")
}

func (r *organizationRepository) FindOne(filter interface{}) (*models.Organization, error) {
	var result models.Organization

	err := r.Collection.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetProjection(bson.D{{Key: "user_id", Value: 0}}),
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *organizationRepository) FindByUser(userId string) (*models.Organization, error) {
	organization, err := r.FindOne(bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return nil, err
	}

	return organization, nil
}

// Get implements Repository.
func (r *organizationRepository) All() (*[]models.Organization, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (r *organizationRepository) Update(organization *models.Organization) (*models.Organization, error) {
	panic("unimplemented")
}

func (r *organizationRepository) GetCollection() *mongo.Collection {
	return r.Collection
}

func NewOrganizationRepository() OrganizationRepository {
	return &organizationRepository{
		Collection: database.DB.Collection("organizations"),
	}
}
