package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	NodeReachEventCollectionName = "node_reach_events"
)

func NodeReachEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.NodeReachEvent, error) {
	var v types.NodeReachEvent
	if err := FindOne(ctx, db.Collection(NodeReachEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func NodeReachEventSave(ctx context.Context, db *mongo.Database, v *types.NodeReachEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(NodeReachEventCollectionName), v, opts...)
}

func NodeReachEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.NodeReachEvent, error) {
	var v types.NodeReachEvent
	if err := FindOneAndUpdate(ctx, db.Collection(NodeReachEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}
