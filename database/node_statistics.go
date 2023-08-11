package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NodeStatisticsCollectionName = "node_statistics"
)

func NodeStatisticsFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	var v bson.M
	if err := FindOne(ctx, db.Collection(NodeStatisticsCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return v, nil
}

func NodeStatisticsSave(ctx context.Context, db *mongo.Database, v bson.M, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(NodeStatisticsCollectionName), v, opts...)
}

func NodeStatisticsFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (bson.M, error) {
	var v bson.M
	if err := FindOneAndUpdate(ctx, db.Collection(NodeStatisticsCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return v, nil
}

func NodeStatisticsFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	var v []bson.M
	if err := FindAll(ctx, db.Collection(NodeStatisticsCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func NodeStatisticsDrop(ctx context.Context, db *mongo.Database) error {
	if err := Drop(ctx, db.Collection(NodeStatisticsCollectionName)); err != nil {
		return err
	}

	return nil
}

func NodeStatisticsDeleteMany(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.DeleteOptions) error {
	if err := DeleteMany(ctx, db.Collection(NodeStatisticsCollectionName), filter, opts...); err != nil {
		return err
	}

	return nil
}

func NodeStatisticsAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(NodeStatisticsCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func NodeStatisticsIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	if err := IndexesCreateMany(ctx, db.Collection(NodeStatisticsCollectionName), models, opts...); err != nil {
		return err
	}

	return nil
}
