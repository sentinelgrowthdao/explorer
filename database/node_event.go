package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	NodeEventCollectionName = "node_events"
)

func NodeEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.NodeEvent, error) {
	var v types.NodeEvent
	if err := FindOne(ctx, db.Collection(NodeEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func NodeEventSave(ctx context.Context, db *mongo.Database, v *types.NodeEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(NodeEventCollectionName), v, opts...)
}

func NodeEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.NodeEvent, error) {
	var v types.NodeEvent
	if err := FindOneAndUpdate(ctx, db.Collection(NodeEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func NodeEventFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.NodeEvent, error) {
	var v []*types.NodeEvent
	if err := FindAll(ctx, db.Collection(NodeEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func NodeEventAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(NodeEventCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
