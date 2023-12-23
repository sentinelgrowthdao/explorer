package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	EventCollectionName = "events"
)

func EventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Event, error) {
	var v models.Event
	if err := FindOne(ctx, db.Collection(EventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func EventInsertOne(ctx context.Context, db *mongo.Database, v *models.Event, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(EventCollectionName), v, opts...)
}

func EventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Event, error) {
	var v models.Event
	if err := FindOneAndUpdate(ctx, db.Collection(EventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func EventFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Event, error) {
	var v []*models.Event
	if err := Find(ctx, db.Collection(EventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func EventIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(EventCollectionName), models, opts...)
}

func EventAggregateAll(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := AggregateAll(ctx, db.Collection(EventCollectionName), pipeline, &v, opts...); err != nil {
		return nil, err
	}

	return v, nil
}
