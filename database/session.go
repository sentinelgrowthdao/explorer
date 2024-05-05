package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	SessionCollectionName = "sessions"
)

func SessionFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Session, error) {
	var v models.Session
	if err := FindOne(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SessionInsertOne(ctx context.Context, db *mongo.Database, v *models.Session, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SessionCollectionName), v, opts...)
}

func SessionFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Session, error) {
	var v models.Session
	if err := FindOneAndUpdate(ctx, db.Collection(SessionCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SessionFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Session, error) {
	var v []*models.Session
	if err := Find(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SessionIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SessionCollectionName), models, opts...)
}

func SessionCountDocuments(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	return CountDocuments(ctx, db.Collection(SessionCollectionName), filter, opts...)
}

func SessionDistinct(ctx context.Context, db *mongo.Database, fieldName string, filter bson.M, opts ...*options.DistinctOptions) (bson.A, error) {
	return Distinct(ctx, db.Collection(SessionCollectionName), fieldName, filter, opts...)
}

func SessionAggregateAll(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := AggregateAll(ctx, db.Collection(SessionCollectionName), pipeline, &v, opts...); err != nil {
		return nil, err
	}

	return v, nil
}
