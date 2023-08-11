package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	DepositEventCollectionName = "deposit_events"
)

func DepositEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.DepositEvent, error) {
	var v types.DepositEvent
	if err := FindOne(ctx, db.Collection(DepositEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func DepositEventSave(ctx context.Context, db *mongo.Database, v *types.DepositEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(DepositEventCollectionName), v, opts...)
}

func DepositEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.DepositEvent, error) {
	var v types.DepositEvent
	if err := FindOneAndUpdate(ctx, db.Collection(DepositEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func DepositEventFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.DepositEvent, error) {
	var v []*types.DepositEvent
	if err := FindAll(ctx, db.Collection(DepositEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func DepositEventsAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(DepositEventCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
