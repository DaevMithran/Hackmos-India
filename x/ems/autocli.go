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
                    Use:       "get [organizer]",
                    Short:     "Get event by organizer",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "organizer"},
                    },
                },
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
                    RpcMethod: "MsgCreateEvent",
                    Use:       "create [name]",
                    Short:     "Create an event",
                    PositionalArgs: []*autocliv1.PositionalArgDescriptor{
                        {ProtoField: "name"},
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
