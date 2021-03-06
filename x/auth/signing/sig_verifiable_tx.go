package signing

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/tx/signing"
)

// SigVerifiableTx defines a Tx interface for all signature verification decorators
type SigVerifiableTx interface {
	types.Tx
	GetSigners() []types.AccAddress
	GetPubKeys() []crypto.PubKey // If signer already has pubkey in context, this list will have nil in its place
	GetSignatures() [][]byte
	GetSignaturesV2() ([]signing.SignatureV2, error)
}

// Tx defines a transaction interface that supports all standard message, signature
// fee, memo, and auxiliary interfaces.
type Tx interface {
	SigVerifiableTx

	types.TxWithMemo
	types.FeeTx
	types.TxWithTimeoutHeight
}
