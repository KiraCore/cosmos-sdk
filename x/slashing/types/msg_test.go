package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"cosmos-sdk/MsgUnjail","value":{"address":"kiravaloper1v93xxeq7m27kg"}}`,
		string(bytes),
	)
}
