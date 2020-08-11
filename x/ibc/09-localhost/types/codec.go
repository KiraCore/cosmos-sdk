package types

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	codectypes "github.com/KiraCore/cosmos-sdk/codec/types"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	clientexported "github.com/KiraCore/cosmos-sdk/x/ibc/02-client/exported"
)

// REMOVE: once simapp uses proto
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ClientState{}, "ibc/client/localhost/ClientState", nil)
}

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateClient{},
	)
	registry.RegisterImplementations(
		(*clientexported.ClientState)(nil),
		&ClientState{},
	)
}

var (
	// SubModuleCdc references the global x/ibc/09-localhost module codec.
	// The actual codec used for serialization should be provided to x/ibc/09-localhost and
	// defined at the application level.
	SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
