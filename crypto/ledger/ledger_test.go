package ledger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	cryptoAmino "github.com/KiraCore/cosmos-sdk/crypto/codec"
	"github.com/KiraCore/cosmos-sdk/crypto/hd"
	"github.com/KiraCore/cosmos-sdk/testutil"
	sdk "github.com/KiraCore/cosmos-sdk/types"
)

func TestErrorHandling(t *testing.T) {
	// first, try to generate a key, must return an error
	// (no panic)
	path := *hd.NewParams(44, 555, 0, false, 0)
	_, err := NewPrivKeySecp256k1Unsafe(path)
	require.Error(t, err)
}

func TestPublicKeyUnsafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeySecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae98721034fef9cd7c4c63588d3b03feb5281b9d232cba34d6f3d71aee59211ffbfe1fe87",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "kirapub1addwnpepqd87l8xhcnrrtzxnkql7k55ph8fr9jarf4hn6udwukfprlalu8lgwlj7566",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, "kira1w34k53py5v5xyluazqpq65agyajavep2nx4wm6",
		addr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)
}

func TestPublicKeyUnsafeHDPath(t *testing.T) {
	expectedAnswers := []string{
		"kirapub1addwnpepqd87l8xhcnrrtzxnkql7k55ph8fr9jarf4hn6udwukfprlalu8lgwlj7566",
		"kirapub1addwnpepqfsdqjr68h7wjg5wacksmqaypasnra232fkgu5sxdlnlu8j22ztxv0w3map",
		"kirapub1addwnpepqw3xwqun6q43vtgw6p4qspq7srvxhcmvq4jrx5j5ma6xy3r7k6dtxtdtpkc",
		"kirapub1addwnpepqvez9lrp09g8w7gkv42y4yr5p6826cu28ydrhrujv862yf4njmqyyzdgh49",
		"kirapub1addwnpepq06hw3enfrtmq8n67teytcmtnrgcr0yntmyt25kdukfjkerdc7lqgpy7w9t",
		"kirapub1addwnpepqg3trf2gd0s2940nckrxherwqhgmm6xd5h4pcnrh4x7y35h6yafmc3cfk59",
		"kirapub1addwnpepqdm6rjpx6wsref8wjn7ym6ntejet430j4szpngfgc20caz83lu545uj33s5",
		"kirapub1addwnpepqvdhtjzy2wf44dm03jxsketxc07vzqwvt3vawqqtljgsr9s7jvydjt98v7m",
		"kirapub1addwnpepqwystfpyxwcava7v3t7ndps5xzu6s553wxcxzmmnxevlzvwrlqpzz2tf9fs",
		"kirapub1addwnpepqw970u6gjqkccg9u3rfj99857wupj2z9fqfzy2w7e5dd7xn7kzzgkcw9pgk",
	}

	const numIters = 10

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, err := NewPrivKeySecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)
		require.NotNil(t, priv)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedAnswers[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func TestPublicKeySafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, addr, err := NewPrivKeySecp256k1(path, "cosmos")

	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Nil(t, ShowAddress(path, priv.PubKey(), sdk.GetConfig().GetBech32AccountAddrPrefix()))

	require.Equal(t, "eb5ae98721034fef9cd7c4c63588d3b03feb5281b9d232cba34d6f3d71aee59211ffbfe1fe87",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "kirapub1addwnpepqd87l8xhcnrrtzxnkql7k55ph8fr9jarf4hn6udwukfprlalu8lgwlj7566",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	require.Equal(t, "kira1w34k53py5v5xyluazqpq65agyajavep2nx4wm6",
		addr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, addr, addr2)
}

func TestPublicKeyHDPath(t *testing.T) {
	expectedPubKeys := []string{
		"kirapub1addwnpepqd87l8xhcnrrtzxnkql7k55ph8fr9jarf4hn6udwukfprlalu8lgwlj7566",
		"kirapub1addwnpepqfsdqjr68h7wjg5wacksmqaypasnra232fkgu5sxdlnlu8j22ztxv0w3map",
		"kirapub1addwnpepqw3xwqun6q43vtgw6p4qspq7srvxhcmvq4jrx5j5ma6xy3r7k6dtxtdtpkc",
		"kirapub1addwnpepqvez9lrp09g8w7gkv42y4yr5p6826cu28ydrhrujv862yf4njmqyyzdgh49",
		"kirapub1addwnpepq06hw3enfrtmq8n67teytcmtnrgcr0yntmyt25kdukfjkerdc7lqgpy7w9t",
		"kirapub1addwnpepqg3trf2gd0s2940nckrxherwqhgmm6xd5h4pcnrh4x7y35h6yafmc3cfk59",
		"kirapub1addwnpepqdm6rjpx6wsref8wjn7ym6ntejet430j4szpngfgc20caz83lu545uj33s5",
		"kirapub1addwnpepqvdhtjzy2wf44dm03jxsketxc07vzqwvt3vawqqtljgsr9s7jvydjt98v7m",
		"kirapub1addwnpepqwystfpyxwcava7v3t7ndps5xzu6s553wxcxzmmnxevlzvwrlqpzz2tf9fs",
		"kirapub1addwnpepqw970u6gjqkccg9u3rfj99857wupj2z9fqfzy2w7e5dd7xn7kzzgkcw9pgk",
	}

	expectedAddrs := []string{
		"kira1w34k53py5v5xyluazqpq65agyajavep2nx4wm6",
		"kira19ewxwemt6uahejvwf44u7dh6tq859tkyujfetd",
		"kira1a07dzdjgjsntxpp75zg7cgatgq0udh3pgzjg6u",
		"kira1qvw52lmn9gpvem8welghrkc52m3zczyh00c77n",
		"kira17m78ka80fqkkw2c4ww0v4xm5nsu2drgr05d4j8",
		"kira1ferh9ll9c452d2p8k2v7heq084guygkn97k0y5",
		"kira10vf3sxmjg96rqq36axcphzfsl74dsntufcc30e",
		"kira1cq83av8cmnar79h0rg7duh9gnr7wkh22hj5889",
		"kira1dszhfrt226jy5rsre7e48vw9tgwe90uenxwhgs",
		"kira1734d7qsylzrdt05muhqqtpd90j8mp4y6ndjexj",
	}

	const numIters = 10

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, addr, err := NewPrivKeySecp256k1(path, "cosmos")
		require.Nil(t, err, "%s", err)
		require.NotNil(t, addr)
		require.NotNil(t, priv)

		addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
		require.Equal(t, addr2, addr)
		require.Equal(t,
			expectedAddrs[i], addr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedPubKeys[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func getFakeTx(accountNumber uint32) []byte {
	tmp := fmt.Sprintf(
		`{"account_number":"%d","chain_id":"1234","fee":{"amount":[{"amount":"150","denom":"atom"}],"gas":"5000"},"memo":"memo","msgs":[[""]],"sequence":"6"}`,
		accountNumber)

	return []byte(tmp)
}

func TestSignaturesHD(t *testing.T) {
	for account := uint32(0); account < 100; account += 30 {
		msg := getFakeTx(account)

		path := *hd.NewFundraiserParams(account, sdk.CoinType, account/5)
		fmt.Printf("Checking signature at %v    ---   PLEASE REVIEW AND ACCEPT IN THE DEVICE\n", path)

		priv, err := NewPrivKeySecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)

		pub := priv.PubKey()
		sig, err := priv.Sign(msg)
		require.Nil(t, err)

		valid := pub.VerifyBytes(msg, sig)
		require.True(t, valid, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)
	}
}

func TestRealDeviceSecp256k1(t *testing.T) {
	msg := getFakeTx(50)
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeySecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)

	pub := priv.PubKey()
	sig, err := priv.Sign(msg)
	require.Nil(t, err)

	valid := pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// now, let's serialize the public key and make sure it still works
	bs := priv.PubKey().Bytes()
	pub2, err := cryptoAmino.PubKeyFromBytes(bs)
	require.Nil(t, err, "%+v", err)

	// make sure we get the same pubkey when we load from disk
	require.Equal(t, pub, pub2)

	// signing with the loaded key should match the original pubkey
	sig, err = priv.Sign(msg)
	require.Nil(t, err)
	valid = pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// make sure pubkeys serialize properly as well
	bs = pub.Bytes()
	bpub, err := cryptoAmino.PubKeyFromBytes(bs)
	require.NoError(t, err)
	require.Equal(t, pub, bpub)
}
