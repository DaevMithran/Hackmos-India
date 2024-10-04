package keeper

import (
	"context"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rollchains/dmhackmoschain/x/ems/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.k.authority, msg.Authority)
	}

	return nil, ms.k.Params.Set(ctx, msg.Params)
}

// MsgCreateEvent implements types.MsgServer.
func (ms msgServer) MsgCreateEvent(ctx context.Context, msg *types.MsgCreateEventParams) (*types.MsgCreateEventResponse, error) {
	organizer, err := ms.k.addressCodec.StringToBytes(msg.Organizer)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid organizer address: %s", err)
	}

	err = ms.k.CreateEvent(ctx, msg.Id, msg.Name, organizer, msg.NftPrice, msg.TokenPrice, msg.TokenSupply, msg.MaxOrganizers)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateEventResponse{}, nil
}

// MsgIssueEventNFT implements types.MsgServer.
func (ms msgServer) MsgIssueEventNFT(ctx context.Context, msg *types.MsgIssueEventNFTParams) (*types.MsgIssueEventNFTResponse, error) {
	isActive, err := ms.k.isActive(ctx, msg.Id)
	if err != nil {
		return nil, err
	}

	if !isActive {
		return nil, fmt.Errorf("event is not active")
	}

	receiver, err := ms.k.addressCodec.StringToBytes(msg.Receiver)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if msg.Nft {
		err = ms.k.MintEventNFT(ctx, receiver, msg.Id)
	} else {
		err = ms.k.IssueEventToken(ctx, receiver, msg.Id)
	}

	if err != nil {
		return nil, err
	} 

	return &types.MsgIssueEventNFTResponse{}, nil
}

// MsgAddEventOrganizer implements types.MsgServer.
func (ms msgServer) MsgAddEventOrganizer(ctx context.Context, msg *types.MsgAddEventOrganizerParams) (*types.MsgAddEventOrganizerResponse, error) {
	isActive, err := ms.k.isActive(ctx, msg.Id)
	if err != nil {
		return nil, err
	}

	if !isActive {
		return nil, fmt.Errorf("event is not active")
	}

	organizer, err := ms.k.addressCodec.StringToBytes(msg.Organizer)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid organizer address: %s", err)
	}

	member, err := ms.k.addressCodec.StringToBytes(msg.Member)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", err)
	}

	isOrganizer, err := ms.k.isOrganizer(ctx, msg.Id, organizer)
	if err != nil {
		return nil, err
	}

	if !isOrganizer {
		return nil, fmt.Errorf("permission denied")
	}

	err = ms.k.AddOrganizer(ctx, organizer, member, msg.Id)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddEventOrganizerResponse{}, nil
}

// MsgUpdateEventStatus implements types.MsgServer.
func (ms msgServer) MsgUpdateEventStatus(ctx context.Context, msg *types.MsgUpdateEventStatusParams) (*types.MsgUpdateEventStatusResponse, error) {
	organizer, err := ms.k.addressCodec.StringToBytes(msg.Organizer)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid organizer address: %s", err)
	}
	err = ms.k.UpdateEventStatus(ctx, organizer, msg.Id, msg.Active)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateEventStatusResponse{}, nil

}
