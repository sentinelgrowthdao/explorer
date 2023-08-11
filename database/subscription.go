package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	SubscriptionCollectionName = "subscriptions"
)

func SubscriptionFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Subscription, error) {
	var v types.Subscription
	if err := FindOne(ctx, db.Collection(SubscriptionCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionSave(ctx context.Context, db *mongo.Database, v *types.Subscription, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(SubscriptionCollectionName), v, opts...)
}

func SubscriptionFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Subscription, error) {
	var v types.Subscription
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Subscription, error) {
	var v []*types.Subscription
	if err := FindAll(ctx, db.Collection(SubscriptionCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(SubscriptionCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionCountDocuments(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	return CountDocuments(ctx, db.Collection(SubscriptionCollectionName), filter, opts...)
}
