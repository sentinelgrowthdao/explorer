package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	SubscriptionQuotaCollectionName = "subscription_quotas"
)

func SubscriptionQuotaFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.SubscriptionQuota, error) {
	var v types.SubscriptionQuota
	if err := FindOne(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionQuotaSave(ctx context.Context, db *mongo.Database, v *types.SubscriptionQuota, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(SubscriptionQuotaCollectionName), v, opts...)
}

func SubscriptionQuotaFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.SubscriptionQuota, error) {
	var v types.SubscriptionQuota
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionQuotaFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.SubscriptionQuota, error) {
	var v []*types.SubscriptionQuota
	if err := FindAll(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
