package provider

import (
	providertypes "github.com/sentinel-official/hub/x/provider/types"
)

type (
	MsgRegisterRequest struct {
		From        string `json:"from,omitempty" bson:"from"`
		Name        string `json:"name,omitempty" bson:"name"`
		Identity    string `json:"identity,omitempty" bson:"identity"`
		Website     string `json:"website,omitempty" bson:"website"`
		Description string `json:"description,omitempty" bson:"description"`
	}
	MsgUpdateRequest struct {
		From        string `json:"from,omitempty" bson:"from"`
		Name        string `json:"name,omitempty" bson:"name"`
		Identity    string `json:"identity,omitempty" bson:"identity"`
		Website     string `json:"website,omitempty" bson:"website"`
		Description string `json:"description,omitempty" bson:"description"`
	}
)

func NewMsgRegisterRequestFromRaw(v *providertypes.MsgRegisterRequest) *MsgRegisterRequest {
	return &MsgRegisterRequest{
		From:        v.From,
		Name:        v.Name,
		Identity:    v.Identity,
		Website:     v.Website,
		Description: v.Description,
	}
}

func NewMsgMsgUpdateRequestFromRaw(v *providertypes.MsgUpdateRequest) *MsgUpdateRequest {
	return &MsgUpdateRequest{
		From:        v.From,
		Name:        v.Name,
		Identity:    v.Identity,
		Website:     v.Website,
		Description: v.Description,
	}
}
