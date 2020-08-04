package types_test

import (
	"testing"

	cryptoAmino "github.com/KiraCore/cosmos-sdk/crypto/codec"

	"github.com/KiraCore/cosmos-sdk/testutil/testdata"

	"github.com/KiraCore/cosmos-sdk/client/testutil"

	"github.com/stretchr/testify/suite"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/auth/types"
)

func testCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := types.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
