package v040_test

import (
	"testing"

	"github.com/KiraCore/cosmos-sdk/codec"
	cryptocodec "github.com/KiraCore/cosmos-sdk/crypto/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	v038auth "github.com/KiraCore/cosmos-sdk/x/auth/legacy/v0_38"
	v039auth "github.com/KiraCore/cosmos-sdk/x/auth/legacy/v0_39"
	v038bank "github.com/KiraCore/cosmos-sdk/x/bank/legacy/v0_38"
	v040bank "github.com/KiraCore/cosmos-sdk/x/bank/legacy/v0_40"

	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	v040Codec := codec.New()
	cryptocodec.RegisterCrypto(v040Codec)
	v039auth.RegisterCodec(v040Codec)

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))
	addr1, err := sdk.AccAddressFromBech32("kira1xxkueklal9vejv9unqu80w9vptyepfa9yw86s3")
	require.NoError(t, err)
	acc1 := v038auth.NewBaseAccount(addr1, coins, nil, 1, 0)

	addr2, err := sdk.AccAddressFromBech32("kira15v50ymp6n5dn73erkqtmq0u8adpl8d3uzrqhlc")
	require.NoError(t, err)
	vaac := v038auth.NewContinuousVestingAccountRaw(
		v038auth.NewBaseVestingAccount(
			v038auth.NewBaseAccount(addr2, coins, nil, 1, 0), coins, nil, nil, 3160620846,
		),
		1580309972,
	)

	bankGenState := v038bank.GenesisState{
		SendEnabled: true,
	}
	authGenState := v039auth.GenesisState{
		Accounts: v038auth.GenesisAccounts{acc1, vaac},
	}

	migrated := v040bank.Migrate(bankGenState, authGenState)
	expected := `{
  "send_enabled": true,
  "balances": [
    {
      "address": "kira1xxkueklal9vejv9unqu80w9vptyepfa9yw86s3",
      "coins": [
        {
          "denom": "stake",
          "amount": "50"
        }
      ]
    },
    {
      "address": "kira15v50ymp6n5dn73erkqtmq0u8adpl8d3uzrqhlc",
      "coins": [
        {
          "denom": "stake",
          "amount": "50"
        }
      ]
    }
  ]
}`

	bz, err := v040Codec.MarshalJSONIndent(migrated, "", "  ")
	require.NoError(t, err)
	require.Equal(t, expected, string(bz))
}
