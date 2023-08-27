package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	SyncStatusCollectionName = "sync_statuses"
)

func SyncStatusFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.SyncStatus, error) {
	var v models.SyncStatus
	if err := FindOne(ctx, db.Collection(SyncStatusCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SyncStatusInsertOne(ctx context.Context, db *mongo.Database, v *models.SyncStatus, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SyncStatusCollectionName), v, opts...)
}

func SyncStatusFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.SyncStatus, error) {
	var v models.SyncStatus
	if err := FindOneAndUpdate(ctx, db.Collection(SyncStatusCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}
