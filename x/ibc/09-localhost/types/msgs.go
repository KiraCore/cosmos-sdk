package types

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	sdkerrors "github.com/KiraCore/cosmos-sdk/types/errors"
	clientexported "github.com/KiraCore/cosmos-sdk/x/ibc/02-client/exported"
	host "github.com/KiraCore/cosmos-sdk/x/ibc/24-host"
)

// Message types for the IBC client
const (
	TypeMsgCreateClient string = "create_client"
)

var (
	_ clientexported.MsgCreateClient = (*MsgCreateClient)(nil)
)

// NewMsgCreateClient creates a new MsgCreateClient instance
func NewMsgCreateClient(signer sdk.AccAddress) *MsgCreateClient {
	return &MsgCreateClient{
		Signer: signer,
	}
}

// Route implements sdk.Msg
func (msg MsgCreateClient) Route() string {
	return host.RouterKey
}

// Type implements sdk.Msg
func (msg MsgCreateClient) Type() string {
	return TypeMsgCreateClient
}

// ValidateBasic implements sdk.Msg
func (msg MsgCreateClient) ValidateBasic() error {
	if msg.Signer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgCreateClient) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgCreateClient) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// GetClientID implements clientexported.MsgCreateClient
func (msg MsgCreateClient) GetClientID() string {
	return clientexported.ClientTypeLocalHost
}

// GetClientType implements clientexported.MsgCreateClient
func (msg MsgCreateClient) GetClientType() string {
	return clientexported.ClientTypeLocalHost
}

// GetConsensusState implements clientexported.MsgCreateClient
func (msg MsgCreateClient) GetConsensusState() clientexported.ConsensusState {
	return nil
}
