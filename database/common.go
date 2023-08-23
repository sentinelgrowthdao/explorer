package database

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findOneError(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

func findOneAndUpdateError(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

func findError(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

func PrepareClient(ctx context.Context, appName, uri, username, password string) (*mongo.Client, error) {
	var (
		registry = bson.NewRegistryBuilder().
				RegisterTypeMapEntry(bsontype.DateTime, reflect.TypeOf(time.Time{})).
				RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{})).
				Build()
		opts = options.Client().
			SetAppName(appName).
			ApplyURI(uri).
			SetRegistry(registry).
			SetMaxPoolSize(0).
			SetMaxConnecting(0).
			SetMinPoolSize(256)
	)

	if username != "" || password != "" {
		opts = opts.SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			},
		)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func PrepareDatabase(ctx context.Context, appName, uri, username, password, name string) (*mongo.Database, error) {
	c, err := PrepareClient(ctx, appName, uri, username, password)
	if err != nil {
		return nil, err
	}

	db := c.Database(name)
	return db, nil
}

func FindOne(ctx context.Context, c *mongo.Collection, filter bson.M, v interface{}, opts ...*options.FindOneOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "FindOne", time.Since(now))
	}()

	result := c.FindOne(ctx, filter, opts...)
	if result.Err() != nil {
		return result.Err()
	}
	if err := result.Decode(v); err != nil {
		return err
	}

	return nil
}

func Save(ctx context.Context, c *mongo.Collection, v interface{}, opts ...*options.InsertOneOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "Save", time.Since(now))
	}()

	_, err := c.InsertOne(ctx, v, opts...)
	return err
}

func FindOneAndUpdate(ctx context.Context, c *mongo.Collection, filter, update bson.M, v interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "FindOneAndUpdate", time.Since(now))
	}()

	result := c.FindOneAndUpdate(ctx, filter, update, opts...)
	if result.Err() != nil {
		return result.Err()
	}
	if err := result.Decode(v); err != nil {
		return err
	}

	return nil
}

func FindAll(ctx context.Context, c *mongo.Collection, filter bson.M, v interface{}, opts ...*options.FindOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "FindAll", time.Since(now))
	}()

	for _, opt := range opts {
		sort, ok := opt.Sort.(bson.D)
		if ok && len(sort) == 0 {
			opt.SetSort(nil)
		}
	}

	cursor, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, v); err != nil {
		return err
	}

	return nil
}

func Aggregate(ctx context.Context, c *mongo.Collection, pipeline []bson.M, v interface{}, opts ...*options.AggregateOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "Aggregate", time.Since(now))
	}()

	cursor, err := c.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, v); err != nil {
		return err
	}

	return nil
}

func CountDocuments(ctx context.Context, c *mongo.Collection, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "CountDocuments", time.Since(now))
	}()

	return c.CountDocuments(ctx, filter, opts...)
}

func Drop(ctx context.Context, c *mongo.Collection) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "Drop", time.Since(now))
	}()

	return c.Drop(ctx)
}

func DeleteMany(ctx context.Context, c *mongo.Collection, filter bson.M, opts ...*options.DeleteOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "DeleteMany", time.Since(now))
	}()

	_, err := c.DeleteMany(ctx, filter, opts...)
	return err
}

func UpdateMany(ctx context.Context, c *mongo.Collection, filter, update bson.M, opts ...*options.UpdateOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "UpdateMany", time.Since(now))
	}()

	_, err := c.UpdateMany(ctx, filter, update, opts...)
	return err
}

func IndexesCreateMany(ctx context.Context, c *mongo.Collection, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	now := time.Now()
	defer func() {
		log.Println(c.Name(), "IndexesCreateMany", time.Since(now))
	}()

	_, err := c.Indexes().CreateMany(ctx, models, opts...)
	if err != nil {
		return err
	}

	return nil
}
