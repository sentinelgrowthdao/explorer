package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	AllocationCollectionName = "allocations"
)

func AllocationFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Allocation, error) {
	var v types.Allocation
	if err := FindOne(ctx, db.Collection(AllocationCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func AllocationSave(ctx context.Context, db *mongo.Database, v *types.Allocation, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(AllocationCollectionName), v, opts...)
}

func AllocationFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Allocation, error) {
	var v types.Allocation
	if err := FindOneAndUpdate(ctx, db.Collection(AllocationCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func AllocationFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Allocation, error) {
	var v []*types.Allocation
	if err := FindAll(ctx, db.Collection(AllocationCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
