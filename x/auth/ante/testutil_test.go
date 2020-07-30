package ante_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	"github.com/KiraCore/cosmos-sdk/simapp"
	"github.com/KiraCore/cosmos-sdk/testutil/testdata"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	txtypes "github.com/KiraCore/cosmos-sdk/types/tx"
	"github.com/KiraCore/cosmos-sdk/types/tx/signing"
	"github.com/KiraCore/cosmos-sdk/x/auth/ante"
	xauthsigning "github.com/KiraCore/cosmos-sdk/x/auth/signing"
	"github.com/KiraCore/cosmos-sdk/x/auth/types"
	authtypes "github.com/KiraCore/cosmos-sdk/x/auth/types"
)

// TestAccount represents an account used in the tests in x/auth/ante.
type TestAccount struct {
	acc  types.AccountI
	priv crypto.PrivKey
}

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	anteHandler sdk.AnteHandler
	ctx         sdk.Context
	clientCtx   client.Context
	txBuilder   client.TxBuilder
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest(isCheckTx bool) {
	suite.app, suite.ctx = createTestApp(isCheckTx)
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Set up TxConfig.
	encodingConfig := simapp.MakeEncodingConfig()
	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	suite.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)

	suite.anteHandler = ante.NewAnteHandler(suite.app.AccountKeeper, suite.app.BankKeeper, ante.DefaultSigVerificationGasConsumer, encodingConfig.TxConfig.SignModeHandler())
}

// CreateTestAccounts creates `numAccs` accounts, and return all relevant
// information about them including their private keys.
func (suite *AnteTestSuite) CreateTestAccounts(numAccs int) []TestAccount {
	var accounts []TestAccount

	for i := 0; i < numAccs; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		err := acc.SetAccountNumber(uint64(i))
		suite.Require().NoError(err)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		suite.app.BankKeeper.SetBalances(suite.ctx, addr, sdk.Coins{
			sdk.NewInt64Coin("atom", 10000000),
		})

		accounts = append(accounts, TestAccount{acc, priv})
	}

	return accounts
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(privs []crypto.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.SigFeeMemoTx, error) {
	// First round: we gather all the signer infos.
	for _, priv := range privs {
		err := suite.txBuilder.SetSignerInfo(priv.PubKey(), &txtypes.ModeInfo{
			Sum: &txtypes.ModeInfo_Single_{
				Single: &txtypes.ModeInfo_Single{
					Mode: suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				},
			},
		})
		if err != nil {
			return nil, err
		}
	}
	// Second round: all signer infos are set, so each signer can sign.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:         chainID,
			AccountNumber:   accNums[i],
			AccountSequence: accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData, suite.txBuilder, priv, suite.clientCtx.TxConfig)
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

// TestCase represents a test case used in test tables.
type TestCase struct {
	desc     string
	malleate func()
	simulate bool
	expPass  bool
	expErr   error
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) RunTestCase(privs []crypto.PrivKey, msgs []sdk.Msg, feeAmount sdk.Coins, gasLimit uint64, accNums, accSeqs []uint64, chainID string, tc TestCase) {
	suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
		suite.Require().NoError(suite.txBuilder.SetMsgs(msgs...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		// Theoretically speaking, ante handler unit tests should only test
		// ante handlers, but here we sometimes also test the tx creation
		// process.
		tx, txErr := suite.CreateTestTx(privs, accNums, accSeqs, chainID)
		newCtx, anteErr := suite.anteHandler(suite.ctx, tx, tc.simulate)

		if tc.expPass {
			suite.Require().NoError(txErr)
			suite.Require().NoError(anteErr)
			suite.Require().NotNil(newCtx)

			suite.ctx = newCtx
		} else {
			switch {
			case txErr != nil:
				suite.Require().Error(txErr)
				suite.Require().True(errors.Is(txErr, tc.expErr))

			case anteErr != nil:
				suite.Require().Error(anteErr)
				suite.Require().True(errors.Is(anteErr, tc.expErr))

			default:
				suite.Fail("expected one of txErr,anteErr to be an error")
			}
		}
	})
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}
