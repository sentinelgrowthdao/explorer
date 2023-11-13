package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StatisticCollectionName = "statistics"
)

func StatisticFindOne(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOneOptions) (bson.M, error) {
	var v bson.M
	if err := FindOne(ctx, db.Collection(StatisticCollectionName), filter, &v, opts...); err != nil {
		return nil, findOneError(err)
	}

	return v, nil
}

func StatisticInsertOne(ctx context.Context, db *mongo.Database, v bson.M, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return InsertOne(ctx, db.Collection(StatisticCollectionName), v, opts...)
}

func StatisticFindOneAndUpdate(ctx context.Context, db *mongo.Database, filter, update bson.M, opts ...*options.FindOneAndUpdateOptions) (bson.M, error) {
	var v bson.M
	if err := FindOneAndUpdate(ctx, db.Collection(StatisticCollectionName), filter, update, &v, opts...); err != nil {
		return nil, findOneAndUpdateError(err)
	}

	return v, nil
}

func StatisticFind(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Find(ctx, db.Collection(StatisticCollectionName), filter, &v, opts...); err != nil {
		return nil, findError(err)
	}

	return v, nil
}

func StatisticDeleteMany(ctx context.Context, db *mongo.Database, filter bson.M, opts ...*options.DeleteOptions) error {
	_, err := DeleteMany(ctx, db.Collection(StatisticCollectionName), filter, opts...)
	if err != nil {
		return err
	}

	return nil
}

func StatisticIndexesCreateMany(ctx context.Context, db *mongo.Database, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return IndexesCreateMany(ctx, db.Collection(StatisticCollectionName), models, opts...)
}

func StatisticAggregate(ctx context.Context, db *mongo.Database, pipeline []bson.M, opts ...*options.AggregateOptions) ([]bson.M, error) {
	var v []bson.M
	if err := Aggregate(ctx, db.Collection(StatisticCollectionName), pipeline, &v, opts...); err != nil {
		return nil, err
	}

	return v, nil
}
