package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	SessionEventCollectionName = "session_events"
)

func SessionEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.SessionEvent, error) {
	var v types.SessionEvent
	if err := FindOne(ctx, db.Collection(SessionEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SessionEventSave(ctx context.Context, db *mongo.Database, v *types.SessionEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(SessionEventCollectionName), v, opts...)
}

func SessionEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.SessionEvent, error) {
	var v types.SessionEvent
	if err := FindOneAndUpdate(ctx, db.Collection(SessionEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SessionEventFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.SessionEvent, error) {
	var v []*types.SessionEvent
	if err := FindAll(ctx, db.Collection(SessionEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SessionEventAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(SessionEventCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
