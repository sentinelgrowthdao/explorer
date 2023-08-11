package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	commontypes "github.com/sentinel-official/explorer/types/common"
)

const (
	appName = "03-node_info"
)

var (
	httpTimeout time.Duration
	dbAddress   string
	dbName      string
	dbUsername  string
	dbPassword  string
)

func init() {
	flag.DurationVar(&httpTimeout, "http-timeout", 5*time.Second, "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

type (
	NodeInfo struct {
		Address                string                 `json:"address"`
		Bandwidth              *commontypes.Bandwidth `json:"bandwidth"`
		Handshake              *types.NodeHandshake   `json:"handshake"`
		IntervalSetSessions    time.Duration          `json:"interval_set_sessions"`
		IntervalUpdateSessions time.Duration          `json:"interval_update_sessions"`
		IntervalUpdateStatus   time.Duration          `json:"interval_update_status"`
		Location               *types.NodeLocation    `json:"location"`
		Moniker                string                 `json:"moniker"`
		Operator               string                 `json:"operator"`
		Peers                  int                    `json:"peers"`
		Price                  string                 `json:"price"`
		Provider               string                 `json:"provider"`
		QOS                    *types.NodeQOS         `json:"qos"`
		Type                   uint64                 `json:"type"`
		Version                string                 `json:"version"`
	}
)

func fetchNodeInfo(remote string) (v NodeInfo, err error) {
	endpoint, err := url.JoinPath(remote, "status")
	if err != nil {
		return v, err
	}

	var (
		rBody      map[string]interface{}
		httpclient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: httpTimeout,
		}
	)

	resp, err := httpclient.Get(endpoint)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&rBody); err != nil {
		return v, err
	}

	res, err := json.Marshal(rBody["result"])
	if err != nil {
		return v, err
	}

	if err := json.Unmarshal(res, &v); err != nil {
		return v, err
	}

	return v, nil
}

func main() {
	db, err := database.PrepareDatabase(context.Background(), appName, dbAddress, dbUsername, dbPassword, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	collNodes := flock.New(filepath.Join(os.TempDir(), "mongodb.sentinelhub-2.nodes.lock"))
	if err := collNodes.Lock(); err != nil {
		log.Fatalln(err)
	}

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

			filter := bson.M{
				"remote_url": bson.M{
					"$exists": true,
				},
				"status": "STATUS_ACTIVE",
			}
			projection := bson.M{
				"address":    1,
				"remote_url": 1,
			}

			dNodes, err := database.NodeFindAll(sctx, db, filter, options.Find().SetProjection(projection))
			if err != nil {
				return err
			}

			var (
				group     = errgroup.Group{}
				timestamp = time.Now()
			)

			for nIndex := 0; nIndex < len(dNodes); nIndex++ {
				index := nIndex
				group.Go(func() error {
					dNodeReachEvent := &types.NodeReachEvent{
						Address:      dNodes[index].Address,
						ErrorMessage: "",
						Timestamp:    timestamp,
					}

					fNode, err := fetchNodeInfo(dNodes[index].RemoteURL)
					if err != nil {
						dNodeReachEvent.ErrorMessage = err.Error()
						log.Println(dNodes[index].Address, dNodes[index].RemoteURL, err)
					}

					if err := database.NodeReachEventSave(sctx, db, dNodeReachEvent); err != nil {
						return err
					}

					filter := bson.M{
						"address": dNodes[index].Address,
					}
					update := bson.M{
						"$set": bson.M{
							"bandwidth":                fNode.Bandwidth,
							"handshake":                fNode.Handshake,
							"interval_set_sessions":    fNode.IntervalSetSessions,
							"interval_update_sessions": fNode.IntervalUpdateStatus,
							"interval_update_status":   fNode.IntervalUpdateStatus,
							"location":                 fNode.Location,
							"moniker":                  fNode.Moniker,
							"peers":                    fNode.Peers,
							"qos":                      fNode.QOS,
							"type":                     fNode.Type,
							"version":                  fNode.Version,
							"reach_status": bson.M{
								"error_message": dNodeReachEvent.ErrorMessage,
								"timestamp":     dNodeReachEvent.Timestamp,
							},
						},
					}
					projection := bson.M{
						"_id": 1,
					}

					_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection))
					if err != nil {
						return err
					}

					return nil
				})
			}

			if err := group.Wait(); err != nil {
				return err
			}

			abort = false
			return sctx.CommitTransaction(sctx)
		},
	)
	log.Println("Duration", time.Since(now))
	if err != nil {
		log.Fatalln(err)
	}

	if err := collNodes.Unlock(); err != nil {
		log.Fatalln(err)
	}
}
