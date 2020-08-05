package types

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	cdctypes "github.com/KiraCore/cosmos-sdk/codec/types"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	clientexported "github.com/KiraCore/cosmos-sdk/x/ibc/02-client/exported"
)

// RegisterCodec registers client state on the provided Amino codec. This type is used for
// Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ClientState{}, "ibc/client/localhost/ClientState", nil)
}

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
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
	amino = codec.New()

	// SubModuleCdc references the global x/ibc/09-localhost module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc/09-localhost and
	// defined at the application level.
	SubModuleCdc = codec.NewHybridCodec(amino, cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
