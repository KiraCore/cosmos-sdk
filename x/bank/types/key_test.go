package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/bank/types"
)

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func TestAddressFromBalancesStore(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira1n88uc38xhjgxzw9nwre4ep2c8ga4fjxcdvs4jn")
	require.NoError(t, err)

	key := cloneAppend(addr.Bytes(), []byte("stake"))
	res := types.AddressFromBalancesStore(key)
	require.Equal(t, res, addr)
}
