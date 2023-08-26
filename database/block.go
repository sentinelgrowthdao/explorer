package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	BlockCollectionName = "blocks"
)

func BlockFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Block, error) {
	var v types.Block
	if err := FindOne(ctx, db.Collection(BlockCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func BlockInsertOne(ctx context.Context, db *mongo.Database, v *types.Block, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(BlockCollectionName), v, opts...)
}

func BlockFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Block, error) {
	var v types.Block
	if err := FindOneAndUpdate(ctx, db.Collection(BlockCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func BlockFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Block, error) {
	var v []*types.Block
	if err := Find(ctx, db.Collection(BlockCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
