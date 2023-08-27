package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	ProviderCollectionName = "providers"
)

func ProviderFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.Provider, error) {
	var v models.Provider
	if err := FindOne(ctx, db.Collection(ProviderCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func ProviderInsertOne(ctx context.Context, db *mongo.Database, v *models.Provider, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(ProviderCollectionName), v, opts...)
}

func ProviderFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.Provider, error) {
	var v models.Provider
	if err := FindOneAndUpdate(ctx, db.Collection(ProviderCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func ProviderFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.Provider, error) {
	var v []*models.Provider
	if err := Find(ctx, db.Collection(ProviderCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
