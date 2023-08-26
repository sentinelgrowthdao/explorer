package utils

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PrepareClient(ctx context.Context, appName, username, password, uri string) (*mongo.Client, error) {
	registry := bson.NewRegistry()
	registry.RegisterTypeMapEntry(bson.TypeDateTime, reflect.TypeOf(time.Time{}))
	registry.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, reflect.TypeOf(bson.M{}))

	opts := options.Client().
		SetAppName(appName).
		ApplyURI(uri).
		SetRegistry(registry).
		SetMaxPoolSize(0).
		SetMaxConnecting(0).
		SetMinPoolSize(256)

	if username != "" && password != "" {
		opts = opts.SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			},
		)
	}

	return mongo.Connect(ctx, opts)
}

func PrepareDatabase(ctx context.Context, appName, username, password, uri, name string) (*mongo.Database, error) {
	client, err := PrepareClient(ctx, appName, username, password, uri)
	if err != nil {
		return nil, err
	}

	return client.Database(name), nil
}
