package keeper_test

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func (suite *IntegrationTestSuite) TestExportGenesis() {
	app, ctx := suite.app, suite.ctx

	expectedMetadata := suite.getTestMetadata()
	expectedBalances := suite.getTestBalances(suite.T())
	for i := range []int{1, 2} {
		app.BankKeeper.SetDenomMetaData(ctx, expectedMetadata[i])
		err := app.BankKeeper.SetBalances(ctx, expectedBalances[i].Address, expectedBalances[i].Coins)
		suite.Require().NoError(err)
	}

	totalSupply := types.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("test", 400000000)))
	app.BankKeeper.SetSupply(ctx, totalSupply)
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	exportGenesis := app.BankKeeper.ExportGenesis(ctx)

	suite.Require().Len(exportGenesis.Params.SendEnabled, 0)
	suite.Require().Equal(types.DefaultParams().DefaultSendEnabled, exportGenesis.Params.DefaultSendEnabled)
	suite.Require().Equal(totalSupply.GetTotal(), exportGenesis.Supply)
	suite.Require().Equal(expectedBalances, exportGenesis.Balances)
	suite.Require().Equal(expectedMetadata, exportGenesis.DenomMetadata)
}

func (suite *IntegrationTestSuite) getTestBalances(t *testing.T) []types.Balance {
	addr2, err := sdk.AccAddressFromBech32("kira1f9xjhxm0plzrh9cskf4qee4pc2xwp0n0ymsxkz")
	require.NoError(t, err)
	addr1, err := sdk.AccAddressFromBech32("kira1fl48vsnmsdzcv85q5d2q4z5ajdha8yu395rpc6")
	require.NoError(t, err)
	return []types.Balance{
		{addr2, sdk.Coins{sdk.NewInt64Coin("testcoin1", 32), sdk.NewInt64Coin("testcoin2", 34)}},
		{addr1, sdk.Coins{sdk.NewInt64Coin("testcoin3", 10)}},
	}

}
