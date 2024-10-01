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
                    Use:       "create <id> <name>",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
						{ProtoField: "name"},
                    },
                },
				{
                    RpcMethod: "MsgIssueEventNFT",
                    Use:       "issue <id> <receiver> <nft>",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "id"},
						{ProtoField: "receiver"},
						{ProtoField: "nft"},
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
