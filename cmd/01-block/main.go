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
	"github.com/sentinel-official/explorer/querier"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "01-block"
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
	flag.Int64Var(&toHeight, "to-height", 1_272_000, "")
	flag.StringVar(&rpcAddress, "rpc-address", "http://127.0.0.1:26657", "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func run(db *mongo.Database, q *querier.Querier, height int64) (operations []types.DatabaseOperation, err error) {
	qBlock, err := q.QueryBlock(context.TODO(), height)
	if err != nil {
		return nil, err
	}

	qBlockResults, err := q.QueryBlockResults(context.TODO(), height)
	if err != nil {
		return nil, err
	}

	operations = append(operations, func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"height": height - 1,
		}
		update := bson.M{
			"$set": bson.M{
				"round":        qBlock.Block.LastCommit.Round,
				"signatures":   types.NewCommitSignatures(qBlock.Block.LastCommit.Signatures),
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

	dBlock := types.NewBlock(qBlock.Block).
		WithBlockID(&qBlock.BlockID).
		WithDuration(qBlock.Block.Time.Sub(prevBlockTime)).
		WithBeginBlockEvents(qBlockResults.BeginBlockEvents).
		WithEndBlockEvents(qBlockResults.EndBlockEvents).
		WithBlockValidatorUpdates(qBlockResults.ValidatorUpdates).
		WithBlockConsensusParams(qBlockResults.ConsensusParamUpdates)
	operations = append(operations, func(ctx mongo.SessionContext) error {
		if _, err := database.BlockInsertOne(ctx, db, dBlock); err != nil {
			return err
		}

		return nil
	})

	log.Println("TxsLen", len(qBlock.Block.Txs))
	for tIndex := 0; tIndex < len(qBlock.Block.Txs); tIndex++ {
		dTx := types.NewTx(qBlock.Block.Txs[tIndex]).
			WithHeight(qBlock.Block.Height).
			WithIndex(tIndex).
			WithResult(qBlockResults.TxsResults[tIndex])
		operations = append(operations, func(ctx mongo.SessionContext) error {
			if _, err := database.TxInsertOne(ctx, db, dTx); err != nil {
				return err
			}

			return nil
		})
	}

	return operations, nil
}

func main() {
	encCfg := hub.MakeEncodingConfig()

	q, err := querier.NewQuerier(&encCfg, rpcAddress, "/websocket")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.TODO(), db, filter)
	if err != nil {
		log.Fatalln(err)
	}
	if dSyncStatus == nil {
		dSyncStatus = &types.SyncStatus{
			AppName:   appName,
			Height:    fromHeight - 1,
			Timestamp: time.Time{},
		}
	}

	height := dSyncStatus.Height + 1
	for height < toHeight {
		now := time.Now()
		log.Println("Height", height)

		operations, err := run(db, q, height)
		if err != nil {
			log.Fatalln(err)
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

				log.Println("OperationsLen", len(operations))
				for i := 0; i < len(operations); i++ {
					if err := operations[i](ctx); err != nil {
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
			log.Fatalln(err)
		}
	}
}
