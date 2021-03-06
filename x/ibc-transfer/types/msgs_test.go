package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

// define constants used for testing
const (
	validPort        = "testportid"
	invalidPort      = "(invalidport1)"
	invalidShortPort = "p"
	invalidLongPort  = "invalidlongportinvalidlongportinvalidlongportinvalidlongportinvalid"

	validChannel        = "testchannel"
	invalidChannel      = "(invalidchannel1)"
	invalidShortChannel = "invalidch"
	invalidLongChannel  = "invalidlongchannelinvalidlongchannelinvalidlongchannelinvalidlongchannel"
)

var (
	addr1     = sdk.AccAddress("testaddr1")
	addr2     = sdk.AccAddress("testaddr2").String()
	emptyAddr sdk.AccAddress

	coin             = sdk.NewCoin("atom", sdk.NewInt(100))
	invalidDenomCoin = sdk.Coin{Denom: "ato-m", Amount: sdk.NewInt(100)}
	negativeCoin     = sdk.Coin{Denom: "atoms", Amount: sdk.NewInt(-100)}
)

// TestMsgTransferRoute tests Route for MsgTransfer
func TestMsgTransferRoute(t *testing.T) {
	msg := NewMsgTransfer(validPort, validChannel, coin, addr1, addr2, 10, 0)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgTransferType tests Type for MsgTransfer
func TestMsgTransferType(t *testing.T) {
	msg := NewMsgTransfer(validPort, validChannel, coin, addr1, addr2, 10, 0)

	require.Equal(t, "transfer", msg.Type())
}

// TestMsgTransferValidation tests ValidateBasic for MsgTransfer
func TestMsgTransferValidation(t *testing.T) {
	testCases := []struct {
		name    string
		msg     *MsgTransfer
		expPass bool
	}{
		{"valid msg", NewMsgTransfer(validPort, validChannel, coin, addr1, addr2, 10, 0), true},
		{"too short port id", NewMsgTransfer(invalidShortPort, validChannel, coin, addr1, addr2, 10, 0), false},
		{"too long port id", NewMsgTransfer(invalidLongPort, validChannel, coin, addr1, addr2, 10, 0), false},
		{"port id contains non-alpha", NewMsgTransfer(invalidPort, validChannel, coin, addr1, addr2, 10, 0), false},
		{"too short channel id", NewMsgTransfer(validPort, invalidShortChannel, coin, addr1, addr2, 10, 0), false},
		{"too long channel id", NewMsgTransfer(validPort, invalidLongChannel, coin, addr1, addr2, 10, 0), false},
		{"channel id contains non-alpha", NewMsgTransfer(validPort, invalidChannel, coin, addr1, addr2, 10, 0), false},
		{"invalid amount", NewMsgTransfer(validPort, validChannel, invalidDenomCoin, addr1, addr2, 10, 0), false},
		{"negative coin", NewMsgTransfer(validPort, validChannel, negativeCoin, addr1, addr2, 10, 0), false},
		{"missing sender address", NewMsgTransfer(validPort, validChannel, coin, emptyAddr, addr2, 10, 0), false},
		{"missing recipient address", NewMsgTransfer(validPort, validChannel, coin, addr1, "", 10, 0), false},
		{"empty coin", NewMsgTransfer(validPort, validChannel, sdk.Coin{}, addr1, addr2, 10, 0), false},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "valid test case %d failed: %v", i, err)
		} else {
			require.Error(t, err, "invalid test case %d passed: %s", i, tc.name)
		}
	}
}

// TestMsgTransferGetSignBytes tests GetSignBytes for MsgTransfer
func TestMsgTransferGetSignBytes(t *testing.T) {
	msg := NewMsgTransfer(validPort, validChannel, coin, addr1, addr2, 110, 10)
	res := msg.GetSignBytes()

	expected := `{"receiver":"kira1w3jhxarpv3j8yvssvrnrt","sender":"kira1w3jhxarpv3j8yvgme08g9","source_channel":"testchannel","source_port":"testportid","timeout_height":"110","timeout_timestamp":"10","token":{"amount":"100","denom":"atom"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgTransferGetSigners tests GetSigners for MsgTransfer
func TestMsgTransferGetSigners(t *testing.T) {
	msg := NewMsgTransfer(validPort, validChannel, coin, addr1, addr2, 10, 0)
	res := msg.GetSigners()

	expected := "[746573746164647231]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
