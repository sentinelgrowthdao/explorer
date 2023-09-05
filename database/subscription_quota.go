package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	SubscriptionQuotaCollectionName = "subscription_quotas"
)

func SubscriptionQuotaFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.SubscriptionQuota, error) {
	var v models.SubscriptionQuota
	if err := FindOne(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionQuotaInsertOne(ctx context.Context, db *mongo.Database, v *models.SubscriptionQuota, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SubscriptionQuotaCollectionName), v, opts...)
}

func SubscriptionQuotaFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.SubscriptionQuota, error) {
	var v models.SubscriptionQuota
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionQuotaFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.SubscriptionQuota, error) {
	var v []*models.SubscriptionQuota
	if err := Find(ctx, db.Collection(SubscriptionQuotaCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionQuotaIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SubscriptionQuotaCollectionName), models, opts...)
}
