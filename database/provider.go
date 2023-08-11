package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	ProviderCollectionName = "providers"
)

func ProviderFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Provider, error) {
	var v types.Provider
	if err := FindOne(ctx, db.Collection(ProviderCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func ProviderSave(ctx context.Context, db *mongo.Database, v *types.Provider, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(ProviderCollectionName), v, opts...)
}

func ProviderFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Provider, error) {
	var v types.Provider
	if err := FindOneAndUpdate(ctx, db.Collection(ProviderCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}
