package client_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/KiraCore/cosmos-sdk/simapp"
	"github.com/KiraCore/cosmos-sdk/testutil"
	"github.com/KiraCore/cosmos-sdk/testutil/testdata"
	authclient "github.com/KiraCore/cosmos-sdk/x/auth/client"

	"github.com/KiraCore/cosmos-sdk/client"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/KiraCore/cosmos-sdk/codec"
	cryptocodec "github.com/KiraCore/cosmos-sdk/crypto/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	authtypes "github.com/KiraCore/cosmos-sdk/x/auth/types"
)

var (
	priv = ed25519.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

func TestParseQueryResponse(t *testing.T) {
	simRes := &sdk.SimulationResponse{
		GasInfo: sdk.GasInfo{GasUsed: 10, GasWanted: 20},
		Result:  &sdk.Result{Data: []byte("tx data"), Log: "log"},
	}

	bz, err := codec.ProtoMarshalJSON(simRes)
	require.NoError(t, err)

	res, err := authclient.ParseQueryResponse(bz)
	require.NoError(t, err)
	require.Equal(t, 10, int(res.GasInfo.GasUsed))
	require.NotNil(t, res.Result)

	res, err = authclient.ParseQueryResponse([]byte("fuzzy"))
	require.Error(t, err)
}

// TODO: remove this and authclient.GetTxEncoder after the proto tx migration is complete
func TestDefaultTxEncoder(t *testing.T) {
	cdc := makeCodec()

	defaultEncoder := authtypes.DefaultTxEncoder(cdc)
	encoder := authclient.GetTxEncoder(cdc)

	compareEncoders(t, defaultEncoder, encoder)
}

func TestReadTxFromFile(t *testing.T) {
	t.Parallel()
	encodingConfig := simapp.MakeEncodingConfig()

	txCfg := encodingConfig.TxConfig
	clientCtx := client.Context{}
	clientCtx = clientCtx.WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
	clientCtx = clientCtx.WithTxConfig(txCfg)

	feeAmount := sdk.Coins{sdk.NewInt64Coin("atom", 150)}
	gasLimit := uint64(50000)
	memo := "foomemo"

	txBuilder := txCfg.NewTxBuilder()
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetMemo(memo)

	// Write it to the file
	encodedTx, err := txCfg.TxJSONEncoder()(txBuilder.GetTx())
	require.NoError(t, err)
	jsonTxFile, cleanup := testutil.WriteToNewTempFile(t, string(encodedTx))
	t.Cleanup(cleanup)

	// Read it back
	decodedTx, err := authclient.ReadTxFromFile(clientCtx, jsonTxFile.Name())
	require.NoError(t, err)
	txBldr, err := txCfg.WrapTxBuilder(decodedTx)
	require.NoError(t, err)
	t.Log(txBuilder.GetTx())
	t.Log(txBldr.GetTx())
	require.Equal(t, txBuilder.GetTx().GetMemo(), txBldr.GetTx().GetMemo())
	require.Equal(t, txBuilder.GetTx().GetFee(), txBldr.GetTx().GetFee())
}

func TestBatchScanner_Scan(t *testing.T) {
	t.Parallel()
	encodingConfig := simapp.MakeEncodingConfig()

	txGen := encodingConfig.TxConfig
	clientCtx := client.Context{}
	clientCtx = clientCtx.WithTxConfig(txGen)

	// generate some tx JSON
	bldr := txGen.NewTxBuilder()
	bldr.SetGasLimit(50000)
	bldr.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("atom", 150)))
	bldr.SetMemo("foomemo")
	txJson, err := txGen.TxJSONEncoder()(bldr.GetTx())
	require.NoError(t, err)

	// use the tx JSON to generate some tx batches (it doesn't matter that we use the same JSON because we don't care about the actual context)
	goodBatchOf3Txs := fmt.Sprintf("%s\n%s\n%s\n", txJson, txJson, txJson)
	malformedBatch := fmt.Sprintf("%s\nmalformed\n%s\n", txJson, txJson)
	batchOf2TxsWithNoNewline := fmt.Sprintf("%s\n%s", txJson, txJson)
	batchWithEmptyLine := fmt.Sprintf("%s\n\n%s", txJson, txJson)

	tests := []struct {
		name               string
		batch              string
		wantScannerError   bool
		wantUnmarshalError bool
		numTxs             int
	}{
		{"good batch", goodBatchOf3Txs, false, false, 3},
		{"malformed", malformedBatch, false, true, 1},
		{"missing trailing newline", batchOf2TxsWithNoNewline, false, false, 2},
		{"empty line", batchWithEmptyLine, false, true, 1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			scanner, i := authclient.NewBatchScanner(clientCtx.TxConfig, strings.NewReader(tt.batch)), 0
			for scanner.Scan() {
				_ = scanner.Tx()
				i++
			}
			require.Equal(t, tt.wantScannerError, scanner.Err() != nil)
			require.Equal(t, tt.wantUnmarshalError, scanner.UnmarshalErr() != nil)
			require.Equal(t, tt.numTxs, i)
		})
	}
}

func compareEncoders(t *testing.T, expected sdk.TxEncoder, actual sdk.TxEncoder) {
	msgs := []sdk.Msg{testdata.NewTestMsg(addr)}
	tx := authtypes.NewStdTx(msgs, authtypes.StdFee{}, []authtypes.StdSignature{}, "")

	defaultEncoderBytes, err := expected(tx)
	require.NoError(t, err)
	encoderBytes, err := actual(tx)
	require.NoError(t, err)
	require.Equal(t, defaultEncoderBytes, encoderBytes)
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	authtypes.RegisterCodec(cdc)
	cdc.RegisterConcrete(testdata.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}
