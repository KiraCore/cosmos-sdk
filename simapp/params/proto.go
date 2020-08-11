// +build !test_amino

package params

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/codec/types"
	"github.com/KiraCore/cosmos-sdk/std"
	"github.com/KiraCore/cosmos-sdk/x/auth/tx"
)

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	amino := codec.New()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewHybridCodec(amino, interfaceRegistry)
	txCfg := tx.NewTxConfig(codec.NewProtoCodec(interfaceRegistry), std.DefaultPublicKeyCodec{}, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}
