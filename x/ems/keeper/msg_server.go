package keeper

import (
	"context"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	// ctx := sdk.UnwrapSDKContext(goCtx)
	err := ms.k.CreateEvent(ctx, msg.Id, msg.Name, sdk.AccAddress(msg.Organizer))
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateEventResponse{}, nil
}

// MsgIssueEventNFT implements types.MsgServer.
func (ms msgServer) MsgIssueEventNFT(ctx context.Context, msg *types.MsgIssueEventNFTParams) (*types.MsgIssueEventNFTResponse, error) {
	event, err := ms.k.GetEvent(ctx, msg.Id);
	if err != nil {
		return nil, err
	}

	isOrganizer := false
	for _, v := range event.Organizers {
		if v == msg.Organizer {
			isOrganizer = true
			break
		}
	}

	if !isOrganizer {
		return nil, fmt.Errorf("permission denied")
	}

    err = ms.k.MintEventNFT(ctx, sdk.AccAddress(msg.Organizer), sdk.AccAddress(msg.Receiver), msg.Id)
	if err != nil {
		return nil, err
	}

	return &types.MsgIssueEventNFTResponse{}, nil
}
