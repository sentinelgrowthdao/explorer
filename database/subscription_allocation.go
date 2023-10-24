package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/models"
)

const (
	SubscriptionAllocationCollectionName = "subscription_allocations"
)

func SubscriptionAllocationFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*models.SubscriptionAllocation, error) {
	var v models.SubscriptionAllocation
	if err := FindOne(ctx, db.Collection(SubscriptionAllocationCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func SubscriptionAllocationInsertOne(ctx context.Context, db *mongo.Database, v *models.SubscriptionAllocation, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(SubscriptionAllocationCollectionName), v, opts...)
}

func SubscriptionAllocationFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*models.SubscriptionAllocation, error) {
	var v models.SubscriptionAllocation
	if err := FindOneAndUpdate(ctx, db.Collection(SubscriptionAllocationCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func SubscriptionAllocationFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*models.SubscriptionAllocation, error) {
	var v []*models.SubscriptionAllocation
	if err := Find(ctx, db.Collection(SubscriptionAllocationCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func SubscriptionAllocationIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(SubscriptionAllocationCollectionName), models, opts...)
}
