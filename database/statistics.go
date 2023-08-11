package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StatisticsCollectionName = "statistics"
)

func StatisticsFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	var v bson.M
	if err := FindOne(ctx, db.Collection(StatisticsCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return v, nil
}

func StatisticsSave(ctx context.Context, db *mongo.Database, v bson.M, opts ...*options.InsertOneOptions) error {
	return Save(ctx, db.Collection(StatisticsCollectionName), v, opts...)
}

func StatisticsFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (bson.M, error) {
	var v bson.M
	if err := FindOneAndUpdate(ctx, db.Collection(StatisticsCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return v, nil
}

func StatisticsFindAll(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	var v []bson.M
	if err := FindAll(ctx, db.Collection(StatisticsCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func StatisticsDrop(ctx context.Context, db *mongo.Database) error {
	if err := Drop(ctx, db.Collection(StatisticsCollectionName)); err != nil {
		return err
	}

	return nil
}

func StatisticsDeleteMany(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.DeleteOptions) error {
	if err := DeleteMany(ctx, db.Collection(StatisticsCollectionName), filter, opts...); err != nil {
		return err
	}

	return nil
}

func StatisticsAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(StatisticsCollectionName), pipeline, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func StatisticsIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	if err := IndexesCreateMany(ctx, db.Collection(StatisticsCollectionName), models, opts...); err != nil {
		return err
	}

	return nil
}
