package models

import (
	"fmt"
	"time"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type CommitSignature struct {
	Flag             string    `json:"flag,omitempty" bson:"flag"`
	Signature        string    `json:"signature,omitempty" bson:"signature"`
	Timestamp        time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	ValidatorAddress string    `json:"validator_address,omitempty" bson:"validator_address"`
}

func NewCommitSignature(v *tmtypes.CommitSig) *CommitSignature {
	return &CommitSignature{
		Flag:             fmt.Sprintf("%v", v.BlockIDFlag),
		Signature:        bytes.HexBytes(v.Signature).String(),
		Timestamp:        v.Timestamp,
		ValidatorAddress: v.ValidatorAddress.String(),
	}
}

type CommitSignatures []*CommitSignature

func NewCommitSignatures(v []tmtypes.CommitSig) CommitSignatures {
	items := make(CommitSignatures, 0, len(v))
	for _, item := range v {
		items = append(items, NewCommitSignature(&item))
	}

	return items
}

type BlockValidatorUpdate struct {
	PubKey string `json:"pub_key,omitempty" bson:"pub_key"`
	Power  int64  `json:"power,omitempty" bson:"power"`
}

func NewBlockValidatorUpdate(v *abcitypes.ValidatorUpdate) *BlockValidatorUpdate {
	pubKey := v.PubKey.GetEd25519()
	if pubKey == nil {
		pubKey = v.PubKey.GetSecp256K1()
	}

	return &BlockValidatorUpdate{
		PubKey: bytes.HexBytes(pubKey).String(),
		Power:  v.Power,
	}
}

type BlockValidatorUpdates []*BlockValidatorUpdate

func NewBlockValidatorUpdates(v []abcitypes.ValidatorUpdate) BlockValidatorUpdates {
	items := make(BlockValidatorUpdates, 0, len(v))
	for _, item := range v {
		items = append(items, NewBlockValidatorUpdate(&item))
	}

	return items
}

type BlockConsensusParams struct {
	BlockMaxBytes           int64    `json:"block_max_bytes,omitempty" bson:"block_max_bytes"`
	BlockMaxGas             int64    `json:"block_max_gas,omitempty" bson:"block_max_gas"`
	EvidenceMaxAgeNumBlocks int64    `json:"evidence_max_age_num_blocks,omitempty" bson:"evidence_max_age_num_blocks"`
	EvidenceMaxAgeDuration  int64    `json:"evidence_max_age_duration,omitempty" bson:"evidence_max_age_duration"`
	EvidenceMaxBytes        int64    `json:"evidence_max_bytes,omitempty" bson:"evidence_max_bytes"`
	ValidatorPubKeyTypes    []string `json:"validator_pub_key_types,omitempty" bson:"validator_pub_key_types"`
	VersionAppVersion       uint64   `json:"version_app_version,omitempty" bson:"version_app_version"`
}

func NewBlockConsensusParams(v *abcitypes.ConsensusParams) *BlockConsensusParams {
	return &BlockConsensusParams{
		BlockMaxBytes:           v.Block.GetMaxBytes(),
		BlockMaxGas:             v.Block.GetMaxBytes(),
		EvidenceMaxAgeNumBlocks: v.Evidence.GetMaxAgeNumBlocks(),
		EvidenceMaxAgeDuration:  v.Evidence.GetMaxAgeDuration().Nanoseconds(),
		EvidenceMaxBytes:        v.Evidence.GetMaxBytes(),
		ValidatorPubKeyTypes:    v.Validator.GetPubKeyTypes(),
		VersionAppVersion:       v.Version.GetAppVersion(),
	}
}

type Block struct {
	AppHash            string                `json:"app_hash,omitempty" bson:"app_hash"`
	BeginBlockEvents   types.Events          `json:"begin_block_events,omitempty" bson:"begin_block_events"`
	ChainID            string                `json:"chain_id,omitempty" bson:"chain_id"`
	CommitHash         string                `json:"commit_hash,omitempty" bson:"commit_hash"`
	ConsensusHash      string                `json:"consensus_hash,omitempty" bson:"consensus_hash"`
	ConsensusParams    *BlockConsensusParams `json:"consensus_params,omitempty" bson:"consensus_params"`
	DataHash           string                `json:"data_hash,omitempty" bson:"data_hash"`
	Duration           int64                 `json:"duration,omitempty" bson:"duration"`
	EndBlockEvents     types.Events          `json:"end_block_events,omitempty" bson:"end_block_events"`
	EvidenceHash       string                `json:"evidence_hash,omitempty" bson:"evidence_hash"`
	Height             int64                 `json:"height,omitempty" bson:"height"`
	ID                 string                `json:"id,omitempty" bson:"id"`
	NextValidatorsHash string                `json:"next_validators_hash,omitempty" bson:"next_validators_hash"`
	NumTxs             int                   `json:"num_txs,omitempty" bson:"num_txs"`
	ProposerAddress    string                `json:"proposer_address,omitempty" bson:"proposer_address"`
	ResultsHash        string                `json:"results_hash,omitempty" bson:"results_hash"`
	Round              int32                 `json:"round,omitempty" bson:"round"`
	Signatures         CommitSignatures      `json:"signatures,omitempty" bson:"signatures"`
	Time               time.Time             `json:"time,omitempty" bson:"time"`
	ValidatorsHash     string                `json:"validators_hash,omitempty" bson:"validators_hash"`
	ValidatorUpdates   BlockValidatorUpdates `json:"validator_updates,omitempty" bson:"validator_updates"`
	Version            string                `json:"version,omitempty" bson:"version"`
}

func NewBlock(v *tmtypes.Block) *Block {
	return &Block{
		AppHash:            v.AppHash.String(),
		BeginBlockEvents:   nil,
		ChainID:            v.ChainID,
		CommitHash:         "",
		ConsensusHash:      v.ConsensusHash.String(),
		ConsensusParams:    nil,
		DataHash:           v.DataHash.String(),
		Duration:           0,
		EndBlockEvents:     nil,
		EvidenceHash:       v.EvidenceHash.String(),
		Height:             v.Height,
		ID:                 "",
		NextValidatorsHash: v.NextValidatorsHash.String(),
		NumTxs:             len(v.Txs),
		ProposerAddress:    v.ProposerAddress.String(),
		ResultsHash:        "",
		Round:              0,
		Signatures:         nil,
		Time:               v.Time,
		ValidatorsHash:     v.ValidatorsHash.String(),
		ValidatorUpdates:   nil,
		Version:            fmt.Sprintf("%d.%d", v.Version.App, v.Version.Block),
	}
}

func (b *Block) String() string {
	return utils.MustMarshalIndentToString(b)
}

func (b *Block) WithBlockID(v *tmtypes.BlockID) *Block   { b.ID = v.String(); return b }
func (b *Block) WithCommitHash(v bytes.HexBytes) *Block  { b.CommitHash = v.String(); return b }
func (b *Block) WithResultsHash(v bytes.HexBytes) *Block { b.ResultsHash = v.String(); return b }
func (b *Block) WithRound(v int32) *Block                { b.Round = v; return b }
func (b *Block) WithDuration(v time.Duration) *Block     { b.Duration = v.Nanoseconds(); return b }

func (b *Block) WithSignatures(v []tmtypes.CommitSig) *Block {
	b.Signatures = NewCommitSignatures(v)
	return b
}

func (b *Block) WithBeginBlockEvents(v []abcitypes.Event) *Block {
	b.BeginBlockEvents = types.NewEventsFromABCIEvents(v)
	return b
}

func (b *Block) WithEndBlockEvents(v []abcitypes.Event) *Block {
	b.EndBlockEvents = types.NewEventsFromABCIEvents(v)
	return b
}

func (b *Block) WithBlockValidatorUpdates(v []abcitypes.ValidatorUpdate) *Block {
	b.ValidatorUpdates = NewBlockValidatorUpdates(v)
	return b
}

func (b *Block) WithBlockConsensusParams(v *abcitypes.ConsensusParams) *Block {
	b.ConsensusParams = NewBlockConsensusParams(v)
	return b
}
