package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/orm/model/ormdb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	apiv1 "github.com/rollchains/dmhackmoschain/api/ems/v1"
	"github.com/rollchains/dmhackmoschain/x/ems/types"
)

type Keeper struct {
	cdc codec.BinaryCodec

	logger log.Logger

	// state management
	Schema collections.Schema
	Params collections.Item[types.Params]
	OrmDB  apiv1.StateStore

	authority string

	EventMapping collections.Map[sdk.AccAddress, string]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	logger log.Logger,
	authority string,
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
		logger: logger,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		OrmDB:  store,

		authority: authority,

		EventMapping: collections.NewMap(sb, collections.NewPrefix(1), "event_mapping", sdk.AccAddressKey, collections.StringValue),
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

func (k Keeper) CreateEvent(ctx context.Context, addr sdk.AccAddress, account string) error {
    has, err := k.EventMapping.Has(ctx, addr)
    if err != nil {
        return err
    }
    if has {
        return fmt.Errorf("account already exists: %s", addr)
    }
    
    err = k.EventMapping.Set(ctx, addr, account)
    if err != nil {
        return err
    }
    return nil
}

func (k Keeper) GetEvent(ctx context.Context, addr sdk.AccAddress) (string, error) {
    acc, err := k.EventMapping.Get(ctx, addr)
    if err != nil {
        return acc, err
    }
    
    return acc, nil
}

func (k Keeper) RemoveEvent(ctx context.Context, addr sdk.AccAddress) error {
    err := k.EventMapping.Remove(ctx, addr)
    if err != nil {
        return err
    }
    return nil
}
