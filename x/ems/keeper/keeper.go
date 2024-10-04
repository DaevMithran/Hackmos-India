package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/orm/model/ormdb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	apiv1 "github.com/rollchains/dmhackmoschain/api/ems/v1"
	"github.com/rollchains/dmhackmoschain/x/ems/types"

	nft "cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
)

type Keeper struct {
	cdc codec.BinaryCodec
	addressCodec address.Codec

	logger log.Logger

	// state management
	Schema collections.Schema
	Params collections.Item[types.Params]
	OrmDB  apiv1.StateStore

	authority string

	EventMapping collections.Map[string, types.Event]

	Nftkeeper nftkeeper.Keeper
	Mintkeeper mintkeeper.Keeper
	BankKeeper bankkeeper.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService storetypes.KVStoreService,
	logger log.Logger,
	authority string,
	nftKeeper nftkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) Keeper {
	logger = logger.With(log.ModuleKey, "x/"+types.ModuleName)

	sb := collections.NewSchemaBuilder(storeService)

	if authority == "" {
		authority = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	}

	db, err := ormdb.NewModuleDB(&types.ORMModuleSchema, ormdb.ModuleDBOptions{KVStoreService: storeService})
	if err != nil {
		panic(err)
	}

	store, err := apiv1.NewStateStore(db)
	if err != nil {
		panic(err)
	}

	k := Keeper{
		cdc:    cdc,
		addressCodec: addressCodec,
		logger: logger,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		OrmDB:  store,

		authority: authority,

		EventMapping: collections.NewMap(sb, collections.NewPrefix(1), "event_mapping", collections.StringKey, codec.CollValue[types.Event](cdc)),

		Nftkeeper: nftKeeper,
		Mintkeeper: mintKeeper,
		BankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) Logger() log.Logger {
	return k.logger
}

// InitGenesis initializes the module's state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	// this line is used by starport scaffolding # genesis/module/init
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return k.Params.Set(ctx, data.Params)
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	// this line is used by starport scaffolding # genesis/module/export

	return &types.GenesisState{
		Params: params,
	}
}

func (k Keeper) CreateEvent(ctx context.Context, id string, name string, account sdk.AccAddress, nftPrice int32, tokenPrice int32, tokenSupply int64, maxOrganizers int32) error {
    has, err := k.EventMapping.Has(ctx, id)
    if err != nil {
        return err
    }
    if has {
        return fmt.Errorf("event already exists: %s", id)
    }

	if nftPrice == 0 {
        nftPrice = 5 // default nftPrice
    }
    if tokenPrice == 0 {
        tokenPrice = 1 // default tokenPrice
    }
    if tokenSupply == 0 {
        tokenSupply = 1000 // default tokenSupply
    }
    if maxOrganizers == 0 {
       maxOrganizers = 10 // default MaxOrganizers
    }
    
    err = k.EventMapping.Set(ctx, id, types.Event {
		Name: name,
		Organizers: []string { account.String() },
		Active: true,
		NftPrice: nftPrice,
		TokenPrice: tokenPrice,
		MaxOrganizers: maxOrganizers,
	})
    if err != nil {
        return err
    }
	// mint tokens
	coins := sdk.NewCoins(sdk.NewCoin(id, math.NewInt(tokenSupply)))
	err = k.Mintkeeper.MintCoins(ctx, coins)
	if err != nil {
		return err
	}

    return nil
}

func (k Keeper) GetEvent(ctx context.Context, id string) (*types.Event, error) {
    event, err := k.EventMapping.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return &event, nil
}

func (k Keeper) RemoveEvent(ctx context.Context, id string) error {
    err := k.EventMapping.Remove(ctx, id)
    if err != nil {
        return err
    }
    return nil
}

func (k Keeper) MintEventNFT(ctx context.Context, receiverAddr sdk.AccAddress, id string) error {
	nftId := receiverAddr.String() + "-" + id

	event, err := k.GetEvent(ctx, id) 
	if err != nil {
		return err
	}

	if !k.Nftkeeper.HasClass(ctx, id) {
		k.Nftkeeper.SaveClass(ctx, nft.Class{ Id: id, Name: event.Name })
	}

	if k.Nftkeeper.HasNFT(ctx, id, nftId) {
		return nil
	}

	coins := sdk.NewCoins(sdk.NewCoin("udmhackmos", math.NewInt(int64(event.NftPrice))))
	// receive payment
	k.BankKeeper.SendCoinsFromAccountToModule(ctx, receiverAddr, types.ModuleName, coins)

	return k.Nftkeeper.Mint(ctx, nft.NFT {
		ClassId: id,
		Id: nftId,
	}, receiverAddr)
}

func (k Keeper) IssueEventToken(ctx context.Context, receiverAddr sdk.AccAddress, id string) error {
	event, err := k.GetEvent(ctx, id) 
	if err != nil {
		return err
	}
	coins := sdk.NewCoins(sdk.NewCoin("udmhackmos", math.NewInt(int64(event.TokenPrice))))
	// receive payment
	k.BankKeeper.SendCoinsFromAccountToModule(ctx, receiverAddr, types.ModuleName, coins)

	// issue token
	eventToken := sdk.NewCoins(sdk.NewCoin(id, math.NewInt(1)))
	return k.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, receiverAddr, eventToken)
}

func (k Keeper) AddOrganizer(ctx context.Context, organizer sdk.AccAddress, member sdk.AccAddress, id string) error {
	isOrg, err := k.isOrganizer(ctx, id, organizer)
	if err != nil {
		return err
	}
	if !isOrg {
		return fmt.Errorf("permission denied")
	}

	event, err := k.GetEvent(ctx, id)
	if err != nil {
		return err
	}

	for _, org := range event.Organizers {
        if org == member.String() {
            return nil
        }
    }
    
    return k.EventMapping.Set(ctx, id, types.Event {
		Name: event.Name,
		Organizers: append(event.Organizers, member.String()),
		Active: event.Active, 
	})
}

func (k Keeper) RemoveOrganizer(ctx context.Context, organizer sdk.AccAddress, organizerToRemove sdk.AccAddress, id string) error {
	isOrg, err := k.isOrganizer(ctx, id, organizer)
	if err != nil {
		return err
	}
	if !isOrg {
		return fmt.Errorf("permission denied")
	}

	isOrg, err = k.isOrganizer(ctx, id, organizerToRemove)
	if err != nil {
		return err
	}

	if !isOrg {
		return nil
	}

	if organizer.String() == organizerToRemove.String() {
		return fmt.Errorf("permission denied")
	}

	event, err := k.GetEvent(ctx, id)
	if err != nil {
		return err
	}

	if event.MaxOrganizers == int32(len(event.Organizers)) {
		return fmt.Errorf("max limit reached for organizers")
	}

	for _, org := range event.Organizers {
        if org == organizerToRemove.String() {
            return nil
        }
    }
    
    return k.EventMapping.Set(ctx, id, types.Event {
		Name: event.Name,
		Organizers: remove(event.Organizers, organizerToRemove.String()) },
	)
}

func remove(slice []string, item string) []string {
    for i, str := range slice {
        if str == item {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}


func (k Keeper) isOrganizer(ctx context.Context, id string, organizer sdk.AccAddress) (bool, error) {
	event, err := k.GetEvent(ctx, id)
	if err != nil {
		return false, err
	}

	for _, v := range event.Organizers {
		if v == organizer.String() {
			return true, nil
		}
	}

	return false, nil
}

func (k Keeper) isActive(ctx context.Context, id string) (bool, error) {
	event, err := k.GetEvent(ctx, id)
	if err != nil {
		return false, err
	}

	return event.Active, nil
}

func (k Keeper) UpdateEventStatus(ctx context.Context, organizer sdk.AccAddress, id string, active bool) error {
	event, err := k.GetEvent(ctx, id)
	if err != nil {
		return err
	}

	if event.Active != active {
		return k.EventMapping.Set(ctx, id, types.Event {
			Name: event.Name,
			Organizers: event.Organizers, 
			Active: active, 
		})
	}

	return nil
}