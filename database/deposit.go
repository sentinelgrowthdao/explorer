package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	DepositCollectionName = "deposits"
)

func DepositFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Deposit, error) {
	var v types.Deposit
	if err := FindOne(ctx, db.Collection(DepositCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func DepositSave(ctx context.Context, db *mongo.Database, v *types.Deposit, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(DepositCollectionName), v, opts...)
}

func DepositFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Deposit, error) {
	var v types.Deposit
	if err := FindOneAndUpdate(ctx, db.Collection(DepositCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func DepositFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Deposit, error) {
	var v []*types.Deposit
	if err := FindAll(ctx, db.Collection(DepositCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
