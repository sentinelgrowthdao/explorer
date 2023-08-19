package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	AllocationEventCollectionName = "allocation_events"
)

func AllocationEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.AllocationEvent, error) {
	var v types.AllocationEvent
	if err := FindOne(ctx, db.Collection(AllocationEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func AllocationEventSave(ctx context.Context, db *mongo.Database, v *types.AllocationEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(AllocationEventCollectionName), v, opts...)
}

func AllocationEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.AllocationEvent, error) {
	var v types.AllocationEvent
	if err := FindOneAndUpdate(ctx, db.Collection(AllocationEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func AllocationEventFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.AllocationEvent, error) {
	var v []*types.AllocationEvent
	if err := FindAll(ctx, db.Collection(AllocationEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
