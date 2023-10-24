package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	SubscriptionPayoutCollectionName = "subscription_payouts"
)

func SubscriptionPayoutFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.SubscriptionPayout, error) {
	var v models.SubscriptionPayout
	if err := FindOne(ctx, db.Collection(SubscriptionPayoutCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionPayoutInsertOne(ctx context.Context, db *mongo.Database, v *models.SubscriptionPayout, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SubscriptionPayoutCollectionName), v, opts...)
}

func SubscriptionPayoutFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.SubscriptionPayout, error) {
	var v models.SubscriptionPayout
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionPayoutCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionPayoutFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.SubscriptionPayout, error) {
	var v []*models.SubscriptionPayout
	if err := Find(ctx, db.Collection(SubscriptionPayoutCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionPayoutIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SubscriptionPayoutCollectionName), models, opts...)
}
