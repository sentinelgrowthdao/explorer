package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	SessionCollectionName = "sessions"
)

func SessionFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Session, error) {
	var v types.Session
	if err := FindOne(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SessionInsertOne(ctx context.Context, db *mongo.Database, v *types.Session, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SessionCollectionName), v, opts...)
}

func SessionFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Session, error) {
	var v types.Session
	if err := FindOneAndUpdate(ctx, db.Collection(SessionCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SessionFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Session, error) {
	var v []*types.Session
	if err := Find(ctx, db.Collection(SessionCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
