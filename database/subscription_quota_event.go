package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	SubscriptionQuotaEventCollectionName = "subscription_quota_events"
)

func SubscriptionQuotaEventFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.SubscriptionQuotaEvent, error) {
	var v types.SubscriptionQuotaEvent
	if err := FindOne(ctx, db.Collection(SubscriptionQuotaEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionQuotaEventSave(ctx context.Context, db *mongo.Database, v *types.SubscriptionQuotaEvent, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(SubscriptionQuotaEventCollectionName), v, opts...)
}

func SubscriptionQuotaEventFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.SubscriptionQuotaEvent, error) {
	var v types.SubscriptionQuotaEvent
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionQuotaEventCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionQuotaEventFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.SubscriptionQuotaEvent, error) {
	var v []*types.SubscriptionQuotaEvent
	if err := FindAll(ctx, db.Collection(SubscriptionQuotaEventCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
