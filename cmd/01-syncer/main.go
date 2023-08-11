package main

import (
	"context"
	"flag"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/querier"
	"github.com/sentinel-official/explorer/types"
)

const (
	appName = "01-syncer"
)

var (
	height     int64
	toHeight   int64
	rpcAddress string
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.Int64Var(&height, "from-height", 9_348_475, "")
	flag.Int64Var(&toHeight, "to-height", math.MaxInt64, "")
	flag.StringVar(&rpcAddress, "rpc-address", "http://127.0.0.1:26657", "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func main() {
	q, err := querier.NewQuerier(rpcAddress, "/websocket")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.PrepareDatabase(context.Background(), appName, dbAddress, dbUsername, dbPassword, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.Background(), db, filter)
	if err != nil {
		log.Fatalln(err)
	}
	if dSyncStatus == nil {
		dSyncStatus = &types.SyncStatus{
			AppName:   appName,
			Height:    height - 1,
			Timestamp: time.Time{},
		}
	}

	height = dSyncStatus.Height + 1

	for height < toHeight {
		log.Println("Height", height)
		now := time.Now()
		err = db.Client().UseSession(
			context.Background(),
			func(sctx mongo.SessionContext) error {
				err = sctx.StartTransaction(
					options.Transaction().
						SetReadConcern(readconcern.Snapshot()).
						SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
				)
				if err != nil {
					return err
				}

				abort := true
				defer func() {
					if abort {
						_ = sctx.AbortTransaction(sctx)
					}
				}()

				qBlockRes, err := q.QueryBlock(sctx, height)
				if err != nil {
					return err
				}

				qBlockResultsRes, err := q.QueryBlockResults(sctx, height)
				if err != nil {
					return err
				}

				filter = bson.M{
					"height": height - 1,
				}
				update := bson.M{
					"$set": bson.M{
						"round":        qBlockRes.Block.LastCommit.Round,
						"signatures":   types.NewCommitSignaturesFromRaw(qBlockRes.Block.LastCommit.Signatures),
						"commit_hash":  qBlockRes.Block.LastCommitHash.String(),
						"results_hash": qBlockRes.Block.LastResultsHash.String(),
					},
				}
				projection := bson.M{
					"time": 1,
				}

				dBlockPrev, err := database.BlockFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection))
				if err != nil {
					return err
				}

				prevBlockTime := qBlockRes.Block.Time
				if dBlockPrev != nil {
					prevBlockTime = dBlockPrev.Time
				}

				dBlock := types.NewBlockFromRaw(qBlockRes.Block).
					WithBlockIDRaw(&qBlockRes.BlockID).
					WithDuration(qBlockRes.Block.Time.Sub(prevBlockTime)).
					WithBeginBlockEventsRaw(qBlockResultsRes.BeginBlockEvents).
					WithEndBlockEventsRaw(qBlockResultsRes.EndBlockEvents).
					WithBlockValidatorUpdatesRaw(qBlockResultsRes.ValidatorUpdates).
					WithBlockConsensusParamsRaw(qBlockResultsRes.ConsensusParamUpdates)

				if err := database.BlockSave(sctx, db, dBlock); err != nil {
					return err
				}

				for tIndex := 0; tIndex < len(qBlockRes.Block.Txs); tIndex++ {
					rTx := types.NewTxFromRaw(qBlockRes.Block.Txs[tIndex]).
						WithHeight(qBlockRes.Block.Height).
						WithIndex(tIndex).
						WithResultRaw(qBlockResultsRes.TxsResults[tIndex])

					if err := database.TxSave(sctx, db, rTx); err != nil {
						return err
					}
				}

				filter = bson.M{
					"app_name": appName,
				}
				update = bson.M{
					"$set": bson.M{
						"height": height,
					},
				}
				projection = bson.M{
					"_id": 1,
				}

				_, err = database.SyncStatusFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
				if err != nil {
					return err
				}

				height++

				abort = false
				return sctx.CommitTransaction(sctx)
			},
		)
		log.Println("Duration", time.Since(now))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
