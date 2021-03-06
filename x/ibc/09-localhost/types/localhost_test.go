package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/simapp"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	clientexported "github.com/KiraCore/cosmos-sdk/x/ibc/02-client/exported"
)

const (
	height = 4
)

type LocalhostTestSuite struct {
	suite.Suite

	cdc   codec.Marshaler
	store sdk.KVStore
}

func (suite *LocalhostTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{Height: 1})
	suite.store = app.IBCKeeper.ClientKeeper.ClientStore(ctx, clientexported.ClientTypeLocalHost)
}

func TestLocalhostTestSuite(t *testing.T) {
	suite.Run(t, new(LocalhostTestSuite))
}
