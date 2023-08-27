package node

import (
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/utils"
)

type MsgRegisterRequest struct {
	From        string
	Name        string
	Identity    string
	Website     string
	Description string `json:"description"`
}

func NewMsgRegisterRequest(v bson.M) (*MsgRegisterRequest, error) {
	return &MsgRegisterRequest{
		From:        v["from"].(string),
		Name:        v["name"].(string),
		Identity:    v["identity"].(string),
		Website:     v["website"].(string),
		Description: v["description"].(string),
	}, nil
}

func (msg *MsgRegisterRequest) ProvAddress() hubtypes.ProvAddress {
	addr := utils.MustAccAddressFromBech32(msg.From)
	return addr.Bytes()
}

type MsgUpdateRequest struct {
	From        string
	Name        string
	Identity    string
	Website     string
	Description string
}

func NewMsgUpdateRequest(v bson.M) (*MsgUpdateRequest, error) {
	return &MsgUpdateRequest{
		From:        v["from"].(string),
		Name:        v["name"].(string),
		Identity:    v["identity"].(string),
		Website:     v["website"].(string),
		Description: v["description"].(string),
	}, nil
}
