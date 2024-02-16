package repository

import (
	"context"
	"log-pulse/models"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChannelRepository interface {
	CreateOne(channel *models.Channel) (*models.Channel, error)
	All(query interface{}, opt *options.FindOptions) (*[]models.Channel, error)
	FindOne(filter interface{}) (*models.Channel, error)
	Update(channel *models.Channel) (*models.Channel, error)
	Delete(filter interface{}, isMany bool) (bool, error)
	Exists(filter interface{}) bool
	GetCollection() *mongo.Collection
}

type channelRepository struct {
	Collection *mongo.Collection
}

func (r *channelRepository) CreateOne(channel *models.Channel) (*models.Channel, error) {
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	channel.ID = ""
	result, err := r.Collection.InsertOne(context.Background(), channel)

	if err != nil {
		return nil, err
	}

	createdChannel, findError := r.FindOne(result.InsertedID)
	if findError != nil {
		return nil, findError
	}

	return createdChannel, nil
}

func (r *channelRepository) FindOne(filter interface{}) (*models.Channel, error) {
	var result models.Channel
	var newFilter interface{}

	filterValue := reflect.ValueOf(filter)
	if filterValue.Type().Kind() == reflect.String || utils.IsObjectID(filter) {
		newFilter = bson.D{{Key: "_id", Value: filter}}
	} else {
		newFilter = filter
	}

	err := r.Collection.FindOne(
		context.Background(),
		newFilter,
		options.FindOne().SetProjection(bson.D{{Key: "user_id", Value: 0}}),
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *channelRepository) All(query interface{}, opt *options.FindOptions) (*[]models.Channel, error) {

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

	var channels []models.Channel = make([]models.Channel, 0)

	if err := cursor.All(context.Background(), &channels); err != nil {
		return nil, err
	}

	return &channels, nil
}

func (r *channelRepository) Delete(filter interface{}, isMany bool) (bool, error) {
	opts := options.Delete().SetHint(bson.D{{Key: "_id", Value: 1}})
	var newFilter interface{}

	var deletedResult *mongo.DeleteResult

	filterValue := reflect.ValueOf(filter)
	if filterValue.Type().Kind() == reflect.String || utils.IsObjectID(filter) {
		newFilter = bson.D{{Key: "_id", Value: filter}}
	} else {
		newFilter = filter
	}

	if isMany == true {
		result, err := r.Collection.DeleteMany(context.Background(), newFilter, opts)
		if err != nil {
			return false, err
		}
		deletedResult = result
	} else {
		result, err := r.Collection.DeleteOne(context.Background(), newFilter, opts)
		if err != nil {
			return false, err
		}

		deletedResult = result
	}

	return deletedResult.DeletedCount > 0, nil
}

func (r *channelRepository) Exists(filter interface{}) bool {
	var channel models.Channel

	err := r.Collection.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}}),
	).Decode(&channel)

	if err != nil {
		return false
	}

	return &channel.ID != nil
}

func (r *channelRepository) Update(model *models.Channel) (*models.Channel, error) {
	panic("unimplemented")
}

func (r *channelRepository) GetCollection() *mongo.Collection {
	return r.Collection
}

func NewChannelRepository() ChannelRepository {
	return &channelRepository{
		Collection: database.DB.Collection("channels"),
	}
}
