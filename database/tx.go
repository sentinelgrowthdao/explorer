package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/types"
)

const (
	TxCollectionName = "txs"
)

func TxFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (*types.Tx, error) {
	var v types.Tx
	if err := FindOne(ctx, db.Collection(TxCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return &v, nil
}

func TxSave(ctx context.Context, db *mongo.Database, v *types.Tx, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(TxCollectionName), v, opts...)
}

func TxFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*types.Tx, error) {
	var v types.Tx
	if err := FindOneAndUpdate(ctx, db.Collection(TxCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return &v, nil
}

func TxFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]*types.Tx, error) {
	var v []*types.Tx
	if err := FindAll(ctx, db.Collection(TxCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}
