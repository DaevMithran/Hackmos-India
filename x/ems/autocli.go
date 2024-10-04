package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	modulev1 "github.com/rollchains/dmhackmoschain/api/ems/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current consensus parameters",
				},
                {
                    RpcMethod: "GetEvent",
                    Use:       "event <id>",
                    Short:     "Get event by id",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
                    },
                },
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
                    RpcMethod: "MsgCreateEvent",
                    Use:       "create <id> <name> [nft_price] [token_price] [token_supply] [max_organizers]",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
						{ProtoField: "name"},
						{ProtoField: "nft_price"},
						{ProtoField: "token_price"},
						{ProtoField: "token_supply"},
						{ProtoField: "max_organizers"},
                    },
                },
				{
                    RpcMethod: "MsgIssueEventNFT",
                    Use:       "issue <id> <nft>",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
						{ProtoField: "nft"},
                    },
                },
				{
                    RpcMethod: "MsgUpdateEventStatus",
                    Use:       "status <id> <active>",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
						{ProtoField: "active"},
                    },
                },
				{
					RpcMethod: "UpdateParams",
					Skip:      false, // set to true if authority gated
				},
			},
		},
	}
}
