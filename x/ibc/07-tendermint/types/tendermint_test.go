package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/simapp"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	ibctmtypes "github.com/KiraCore/cosmos-sdk/x/ibc/07-tendermint/types"
)

const (
	chainID                      = "gaia"
	height                       = 4
	trustingPeriod time.Duration = time.Hour * 24 * 7 * 2
	ubdPeriod      time.Duration = time.Hour * 24 * 7 * 3
	maxClockDrift  time.Duration = time.Second * 10
)

type TendermintTestSuite struct {
	suite.Suite

	ctx      sdk.Context
	aminoCdc *codec.Codec
	cdc      codec.Marshaler
	privVal  tmtypes.PrivValidator
	valSet   *tmtypes.ValidatorSet
	valsHash tmbytes.HexBytes
	header   ibctmtypes.Header
	now      time.Time
}

func (suite *TendermintTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)

	suite.aminoCdc = app.Codec()
	suite.cdc = app.AppCodec()

	suite.now = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	suite.privVal = tmtypes.NewMockPV()

	pubKey, err := suite.privVal.GetPubKey()
	suite.Require().NoError(err)

	val := tmtypes.NewValidator(pubKey, 10)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{val})
	suite.valsHash = suite.valSet.Hash()
	suite.header = ibctmtypes.CreateTestHeader(chainID, height, height-1, suite.now, suite.valSet, suite.valSet, []tmtypes.PrivValidator{suite.privVal})
	suite.ctx = app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, Time: suite.now})
}

func TestTendermintTestSuite(t *testing.T) {
	suite.Run(t, new(TendermintTestSuite))
}
