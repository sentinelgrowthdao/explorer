package types

import (
	"fmt"
	"time"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/sentinel-official/explorer/utils"
)

type CommitSignature struct {
	Flag             string    `json:"flag,omitempty" bson:"flag"`
	ValidatorAddress string    `json:"validator_address,omitempty" bson:"validator_address"`
	Timestamp        time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	Signature        string    `json:"signature,omitempty" bson:"signature"`
}

func NewCommitSignatureFromRaw(v *tmtypes.CommitSig) *CommitSignature {
	return &CommitSignature{
		Flag:             fmt.Sprintf("%v", v.BlockIDFlag),
		ValidatorAddress: v.ValidatorAddress.String(),
		Timestamp:        v.Timestamp,
		Signature:        bytes.HexBytes(v.Signature).String(),
	}
}

type CommitSignatures []*CommitSignature

func NewCommitSignaturesFromRaw(v []tmtypes.CommitSig) CommitSignatures {
	items := make(CommitSignatures, 0, len(v))
	for _, item := range v {
		items = append(items, NewCommitSignatureFromRaw(&item))
	}

	return items
}

type BlockValidatorUpdate struct {
	PubKey string `json:"pub_key,omitempty" bson:"pub_key"`
	Power  int64  `json:"power,omitempty" bson:"power"`
}

func NewBlockValidatorUpdateFromRaw(v *abcitypes.ValidatorUpdate) *BlockValidatorUpdate {
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

func NewBlockValidatorUpdatesFromRaw(v []abcitypes.ValidatorUpdate) BlockValidatorUpdates {
	items := make(BlockValidatorUpdates, 0, len(v))
	for _, item := range v {
		items = append(items, NewBlockValidatorUpdateFromRaw(&item))
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

func NewBlockConsensusParamsFromRaw(v *abcitypes.ConsensusParams) *BlockConsensusParams {
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
	ID                 string                `json:"id,omitempty" bson:"id"`
	Version            string                `json:"version,omitempty" bson:"version"`
	ChainID            string                `json:"chain_id,omitempty" bson:"chain_id"`
	Height             int64                 `json:"height,omitempty" bson:"height"`
	Time               time.Time             `json:"time,omitempty" bson:"time"`
	CommitHash         string                `json:"commit_hash,omitempty" bson:"commit_hash"`
	DataHash           string                `json:"data_hash,omitempty" bson:"data_hash"`
	ValidatorsHash     string                `json:"validators_hash,omitempty" bson:"validators_hash"`
	NextValidatorsHash string                `json:"next_validators_hash,omitempty" bson:"next_validators_hash"`
	ConsensusHash      string                `json:"consensus_hash,omitempty" bson:"consensus_hash"`
	AppHash            string                `json:"app_hash,omitempty" bson:"app_hash"`
	ResultsHash        string                `json:"results_hash,omitempty" bson:"results_hash"`
	EvidenceHash       string                `json:"evidence_hash,omitempty" bson:"evidence_hash"`
	ProposerAddress    string                `json:"proposer_address,omitempty" bson:"proposer_address"`
	NumTxs             int                   `json:"num_txs,omitempty" bson:"num_txs"`
	Round              int32                 `json:"round,omitempty" bson:"round"`
	Duration           int64                 `json:"duration,omitempty" bson:"duration"`
	Signatures         CommitSignatures      `json:"signatures,omitempty" bson:"signatures"`
	BeginBlockEvents   Events                `json:"begin_block_events,omitempty" bson:"begin_block_events"`
	EndBlockEvents     Events                `json:"end_block_events,omitempty" bson:"end_block_events"`
	ValidatorUpdates   BlockValidatorUpdates `json:"validator_updates,omitempty" bson:"validator_updates"`
	ConsensusParams    *BlockConsensusParams `json:"consensus_params,omitempty" bson:"consensus_params"`
}

func NewBlockFromRaw(v *tmtypes.Block) *Block {
	return &Block{
		ID:                 "",
		Version:            fmt.Sprintf("%d.%d", v.Version.App, v.Version.Block),
		ChainID:            v.ChainID,
		Height:             v.Height,
		Time:               v.Time,
		CommitHash:         "",
		DataHash:           v.DataHash.String(),
		ValidatorsHash:     v.ValidatorsHash.String(),
		NextValidatorsHash: v.NextValidatorsHash.String(),
		ConsensusHash:      v.ConsensusHash.String(),
		AppHash:            v.AppHash.String(),
		ResultsHash:        "",
		EvidenceHash:       v.EvidenceHash.String(),
		ProposerAddress:    v.ProposerAddress.String(),
		NumTxs:             len(v.Txs),
		Round:              0,
		Duration:           0,
		Signatures:         nil,
		BeginBlockEvents:   nil,
		EndBlockEvents:     nil,
		ValidatorUpdates:   nil,
		ConsensusParams:    nil,
	}
}

func (b *Block) String() string {
	return utils.MustMarshalIndent(b)
}

func (b *Block) WithBlockIDRaw(v *tmtypes.BlockID) *Block   { b.ID = v.String(); return b }
func (b *Block) WithCommitHashRaw(v bytes.HexBytes) *Block  { b.CommitHash = v.String(); return b }
func (b *Block) WithResultsHashRaw(v bytes.HexBytes) *Block { b.ResultsHash = v.String(); return b }
func (b *Block) WithRound(v int32) *Block                   { b.Round = v; return b }
func (b *Block) WithDuration(v time.Duration) *Block        { b.Duration = v.Nanoseconds(); return b }

func (b *Block) WithSignaturesRaw(v []tmtypes.CommitSig) *Block {
	b.Signatures = NewCommitSignaturesFromRaw(v)
	return b
}

func (b *Block) WithBeginBlockEventsRaw(v []abcitypes.Event) *Block {
	b.BeginBlockEvents = NewEventsFromRaw(v)
	return b
}

func (b *Block) WithEndBlockEventsRaw(v []abcitypes.Event) *Block {
	b.EndBlockEvents = NewEventsFromRaw(v)
	return b
}

func (b *Block) WithBlockValidatorUpdatesRaw(v []abcitypes.ValidatorUpdate) *Block {
	b.ValidatorUpdates = NewBlockValidatorUpdatesFromRaw(v)
	return b
}

func (b *Block) WithBlockConsensusParamsRaw(v *abcitypes.ConsensusParams) *Block {
	b.ConsensusParams = NewBlockConsensusParamsFromRaw(v)
	return b
}
