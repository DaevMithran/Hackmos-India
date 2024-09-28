package keeper

import (
	"context"

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
	err := ms.k.CreateEvent(ctx, sdk.AccAddress(msg.Organizer), msg.Name)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateEventResponse{}, nil
}
