package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/sentinel-official/hub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/querier"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "01-tendermint"
)

var (
	fromHeight int64
	toHeight   int64
	rpcAddress string
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.Int64Var(&fromHeight, "from-height", 901_801, "")
	flag.Int64Var(&toHeight, "to-height", 5_125_000, "")
	flag.StringVar(&rpcAddress, "rpc-address", "http://127.0.0.1:26657", "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "app_name", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err := database.SyncStatusIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "height", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.BlockIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "height", Value: 1},
				bson.E{Key: "result.code", Value: 1},
			},
			Options: options.Index().
				SetPartialFilterExpression(
					bson.M{
						"result.code": 0,
					},
				),
		},
	}

	_, err = database.TxIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
}

func run(db *mongo.Database, q *querier.Querier, height int64) (ops []types.DatabaseOperation, err error) {
	qBlock, err := q.QueryBlock(context.TODO(), height)
	if err != nil {
		return nil, err
	}

	qBlockResults, err := q.QueryBlockResults(context.TODO(), height)
	if err != nil {
		return nil, err
	}

	ops = append(ops, func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"height": height - 1,
		}
		update := bson.M{
			"$set": bson.M{
				"round":        qBlock.Block.LastCommit.Round,
				"signatures":   models.NewCommitSignatures(qBlock.Block.LastCommit.Signatures),
				"commit_hash":  qBlock.Block.LastCommitHash.String(),
				"results_hash": qBlock.Block.LastResultsHash.String(),
			},
		}
		projection := bson.M{
			"time": 1,
		}

		_, err := database.BlockFindOneAndUpdate(ctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection))
		if err != nil {
			return err
		}

		return nil
	})

	filter := bson.M{
		"height": height - 1,
	}

	dBlockPrev, err := database.BlockFindOne(context.TODO(), db, filter)
	if err != nil {
		return nil, err
	}

	prevBlockTime := qBlock.Block.Time
	if dBlockPrev != nil {
		prevBlockTime = dBlockPrev.Time
	}

	dBlock := models.NewBlock(qBlock.Block).
		WithBlockID(&qBlock.BlockID).
		WithDuration(qBlock.Block.Time.Sub(prevBlockTime)).
		WithBeginBlockEvents(qBlockResults.BeginBlockEvents).
		WithEndBlockEvents(qBlockResults.EndBlockEvents).
		WithBlockValidatorUpdates(qBlockResults.ValidatorUpdates).
		WithBlockConsensusParams(qBlockResults.ConsensusParamUpdates)
	ops = append(ops, func(ctx mongo.SessionContext) error {
		if _, err := database.BlockInsertOne(ctx, db, dBlock); err != nil {
			return err
		}

		return nil
	})

	log.Println("TxsLen", len(qBlock.Block.Txs))
	for tIndex := 0; tIndex < len(qBlock.Block.Txs); tIndex++ {
		dTx := models.NewTx(qBlock.Block.Txs[tIndex]).
			WithHeight(qBlock.Block.Height).
			WithIndex(tIndex).
			WithResult(qBlockResults.TxsResults[tIndex]).
			WithTimestamp(qBlock.Block.Time)
		ops = append(ops, func(ctx mongo.SessionContext) error {
			if _, err := database.TxInsertOne(ctx, db, dTx); err != nil {
				return err
			}

			return nil
		})
	}

	return ops, nil
}

func main() {
	encCfg := hub.MakeEncodingConfig()

	q, err := querier.NewQuerier(encCfg.InterfaceRegistry, rpcAddress, "/websocket")
	if err != nil {
		log.Panicln(err)
	}

	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Panicln(err)
	}

	if err := db.Client().Ping(context.TODO(), nil); err != nil {
		log.Panicln(err)
	}

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Panicln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.TODO(), db, filter)
	if err != nil {
		log.Panicln(err)
	}
	if dSyncStatus == nil {
		dSyncStatus = &models.SyncStatus{
			AppName:   appName,
			Height:    fromHeight - 1,
			Timestamp: time.Time{},
		}
	}

	height := dSyncStatus.Height + 1
	for height < toHeight {
		now := time.Now()
		log.Println("Height", height)

		ops, err := run(db, q, height)
		if err != nil {
			log.Panicln(err)
		}

		err = db.Client().UseSession(
			context.TODO(),
			func(ctx mongo.SessionContext) error {
				err := ctx.StartTransaction(
					options.Transaction().
						SetReadConcern(readconcern.Snapshot()).
						SetWriteConcern(writeconcern.Majority()),
				)
				if err != nil {
					return err
				}

				abort := true
				defer func() {
					if abort {
						_ = ctx.AbortTransaction(ctx)
					}
				}()

				log.Println("OperationsLen", len(ops))
				for i := 0; i < len(ops); i++ {
					if err := ops[i](ctx); err != nil {
						return err
					}
				}

				filter := bson.M{
					"app_name": appName,
				}
				update := bson.M{
					"$set": bson.M{
						"height": height,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				_, err = database.SyncStatusFindOneAndUpdate(ctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
				if err != nil {
					return err
				}

				height++

				abort = false
				return ctx.CommitTransaction(ctx)
			},
		)

		log.Println("Duration", time.Since(now))
		log.Println()
		if err != nil {
			log.Panicln(err)
		}
	}
}
