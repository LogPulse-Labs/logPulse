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

type projectRepository struct {
	Collection *mongo.Collection
}

type ProjectRepository interface {
	CreateOne(project *models.Project) (*models.Project, error)
	All(query interface{}, opt *options.FindOptions) (*[]models.Project, error)
	Paginate(query interface{}, options *options.FindOptions) (*[]models.Project, error)
	FindOne(filter interface{}, opt *options.FindOneOptions) (*models.Project, error)
	Update(project *models.Project) (*models.Project, error)
	Delete(ID interface{}) (bool, error)
	Exists(filter interface{}) bool
	GetCollection() *mongo.Collection
}

func (r *projectRepository) CreateOne(project *models.Project) (*models.Project, error) {
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	project.ID = ""
	result, err := r.Collection.InsertOne(context.Background(), project)

	if err != nil {
		return nil, err
	}

	createdProject, findError := r.FindOne(result.InsertedID, nil)
	if findError != nil {
		return nil, findError
	}

	return createdProject, nil
}

func (r *projectRepository) Delete(ID interface{}) (bool, error) {
	opts := options.Delete().SetHint(bson.D{{Key: "_id", Value: 1}})

	result, err := r.Collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: ID}}, opts)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (r *projectRepository) FindOne(filter interface{}, opt *options.FindOneOptions) (*models.Project, error) {
	var result models.Project

	var opts *options.FindOneOptions

	if opt == nil {
		opts = options.FindOne().SetProjection(bson.D{{Key: "user_id", Value: 0}})
	} else {
		opts = opt
	}

	err := r.Collection.FindOne(
		context.Background(),
		bson.D{{Key: "_id", Value: filter}},
		opts,
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *projectRepository) All(query interface{}, opt *options.FindOptions) (*[]models.Project, error) {

	var opts *options.FindOptions

	if opt == nil {
		opts = options.Find().SetProjection(bson.D{{Key: "user_id", Value: 0}})
	} else {
		opts = opt
	}

	opts = opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.Collection.Find(context.Background(), query, opts)

	if err != nil {
		return nil, err
	}

	var projects []models.Project = make([]models.Project, 0)

	if err := cursor.All(context.Background(), &projects); err != nil {
		return nil, err
	}

	return &projects, nil
}

func (r *projectRepository) Paginate(query interface{}, options *options.FindOptions) (*[]models.Project, error) {
	cursor, err := r.Collection.Find(
		context.Background(),
		query,
		options.SetProjection(bson.D{{Key: "user_id", Value: 0}}),
	)

	if err != nil {
		return nil, err
	}

	var projects []models.Project = make([]models.Project, 0)

	if err := cursor.All(context.Background(), &projects); err != nil {
		return nil, err
	}

	return &projects, nil
}

func (r *projectRepository) Exists(filter interface{}) bool {
	var project models.Project

	err := r.Collection.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}}),
	).Decode(&project)

	if err != nil {
		return false
	}

	return &project.ID != nil
}

func (r *projectRepository) Update(model *models.Project) (*models.Project, error) {
	panic("unimplemented")
}

func (r *projectRepository) GetCollection() *mongo.Collection {
	return r.Collection
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{
		Collection: database.DB.Collection("projects"),
	}
}
