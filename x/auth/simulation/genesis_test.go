package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/types/module"
	simtypes "github.com/KiraCore/cosmos-sdk/types/simulation"
	"github.com/KiraCore/cosmos-sdk/x/auth/simulation"
	"github.com/KiraCore/cosmos-sdk/x/auth/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abonormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	cdc := codec.New()
	// Make sure to register cdc.
	// otherwise the test will panic
	types.RegisterCodec(cdc)

	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: 1000,
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var authGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &authGenesis)

	require.Equal(t, uint64(0x8c), authGenesis.Params.GetMaxMemoCharacters())
	require.Equal(t, uint64(0x2b6), authGenesis.Params.GetSigVerifyCostED25519())
	require.Equal(t, uint64(0x1ff), authGenesis.Params.GetSigVerifyCostSecp256k1())
	require.Equal(t, uint64(9), authGenesis.Params.GetTxSigLimit())
	require.Equal(t, uint64(5), authGenesis.Params.GetTxSizeCostPerByte())

	genAccounts, err := types.UnpackAccounts(authGenesis.Accounts)
	require.NoError(t, err)
	require.Len(t, genAccounts, 3)
	require.Equal(t, "kira1ghekyjucln7y67ntx7cf27m9dpuxxemn9hdx5w", genAccounts[2].GetAddress().String())
	require.Equal(t, uint64(0), genAccounts[2].GetAccountNumber())
	require.Equal(t, uint64(0), genAccounts[2].GetSequence())
}