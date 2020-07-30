package v039_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v038auth "github.com/KiraCore/cosmos-sdk/x/auth/legacy/v0_38"
	v039auth "github.com/KiraCore/cosmos-sdk/x/auth/legacy/v0_39"
	v039 "github.com/KiraCore/cosmos-sdk/x/genutil/legacy/v0_39"
	"github.com/KiraCore/cosmos-sdk/x/genutil/types"
)

var genAuthState = []byte(`{
  "params": {
    "max_memo_characters": "10",
    "tx_sig_limit": "10",
    "tx_size_cost_per_byte": "10",
    "sig_verify_cost_ed25519": "10",
    "sig_verify_cost_secp256k1": "10"
  },
  "accounts": [
    {
      "type": "cosmos-sdk/BaseAccount",
      "value": {
        "address": "kira19hz3ee9e3lj9mne4jggj3v8hxjrpre22znuuyf",
        "coins": [
          {
            "denom": "stake",
            "amount": "400000"
          }
        ],
        "public_key": "kirapub1addwnpepqtezq4ajkevh724ls45zp72x70rj8mhszqf5pxcaahazm8trv490s735rgq",
        "account_number": 1,
        "sequence": 1
      }
    },
    {
      "type": "cosmos-sdk/ModuleAccount",
      "value": {
        "address": "kira1fl48vsnmsdzcv85q5d2q4z5ajdha8yu395rpc6",
        "coins": [
          {
            "denom": "stake",
            "amount": "400000000"
          }
        ],
        "public_key": "",
        "account_number": 2,
        "sequence": 4,
        "name": "bonded_tokens_pool",
        "permissions": [
          "burner",
          "staking"
        ]
      }
    },
    {
      "type": "cosmos-sdk/ContinuousVestingAccount",
      "value": {
        "address": "kira1vtzxzyjv506dvhl9pa527xsugf5gez4fr0vww7",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "kirapub1addwnpepqdxrk48q89xlmnzrr5nkssle05tkp73uknevzaavm53c02v26vlyzj53w2z",
        "account_number": 3,
        "sequence": 5,
        "original_vesting": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "delegated_free": [],
        "delegated_vesting": [],
        "end_time": 1596125048,
        "start_time": 1595952248
      }
    },
    {
      "type": "cosmos-sdk/DelayedVestingAccount",
      "value": {
        "address": "kira1prxkcqclweqa0g28p7vmf6z78ghyeckm90hcsz",
        "coins": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "public_key": "kirapub1addwnpepqwewcad349e2yw3weatf8lzfyv5cd6am9jkk4ajach3f568k6gg47r3d8xj",
        "account_number": 4,
        "sequence": 15,
        "original_vesting": [
          {
            "denom": "stake",
            "amount": "10000205"
          }
        ],
        "delegated_free": [],
        "delegated_vesting": [],
        "end_time": 1596125048
      }
    }
  ]
}`)

var expectedGenAuthState = []byte(`{"params":{"max_memo_characters":"10","tx_sig_limit":"10","tx_size_cost_per_byte":"10","sig_verify_cost_ed25519":"10","sig_verify_cost_secp256k1":"10"},"accounts":[{"type":"cosmos-sdk/BaseAccount","value":{"address":"kira19hz3ee9e3lj9mne4jggj3v8hxjrpre22znuuyf","coins":[{"denom":"stake","amount":"400000"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"AvIgV7K2WX8qv4VoIPlG88cj7vAQE0CbHe36LZ1jZUr4"},"account_number":"1","sequence":"1"}},{"type":"cosmos-sdk/ModuleAccount","value":{"address":"kira1fl48vsnmsdzcv85q5d2q4z5ajdha8yu395rpc6","coins":[{"denom":"stake","amount":"400000000"}],"public_key":"","account_number":"2","sequence":"4","name":"bonded_tokens_pool","permissions":["burner","staking"]}},{"type":"cosmos-sdk/ContinuousVestingAccount","value":{"address":"kira1vtzxzyjv506dvhl9pa527xsugf5gez4fr0vww7","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"A0w7VOA5Tf3MQx0naEP5fRdg+jy08sF3rN0jh6mK0z5B"},"account_number":"3","sequence":"5","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048","start_time":"1595952248"}},{"type":"cosmos-sdk/DelayedVestingAccount","value":{"address":"kira1prxkcqclweqa0g28p7vmf6z78ghyeckm90hcsz","coins":[{"denom":"stake","amount":"10000205"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"A7LsdbGpcqI6Ls9Wk/xJIymG67ssrWr2XcXimmj20hFf"},"account_number":"4","sequence":"15","original_vesting":[{"denom":"stake","amount":"10000205"}],"delegated_free":[],"delegated_vesting":[],"end_time":"1596125048"}}]}`)

func TestMigrate(t *testing.T) {
	genesis := types.AppMap{
		v038auth.ModuleName: genAuthState,
	}

	var migrated types.AppMap
	require.NotPanics(t, func() { migrated = v039.Migrate(genesis) })
	require.Equal(t, string(expectedGenAuthState), string(migrated[v039auth.ModuleName]))
}
