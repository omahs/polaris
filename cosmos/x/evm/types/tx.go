// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"

	"pkg.berachain.dev/polaris/eth/common"
	coretypes "pkg.berachain.dev/polaris/eth/core/types"
	"pkg.berachain.dev/polaris/lib/utils"

	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// WrappedEthereumTransaction defines a Cosmos SDK message for Ethereum transactions.
var (
	_ sdk.Msg   = (*WrappedEthereumTransaction)(nil)
	_ sdk.FeeTx = (*WrappedEthereumTransaction)(nil)
)

// NewFromTransaction sets the transaction data from an `coretypes.Transaction`.
func NewFromTransaction(tx *coretypes.Transaction) *WrappedEthereumTransaction {
	bz, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return &WrappedEthereumTransaction{
		Data: bz,
	}
}

// GetSigners returns the address(es) that must sign over the transaction.
func (etr *WrappedEthereumTransaction) GetSigners() []sdk.AccAddress {
	sender, err := etr.GetSender()
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sdk.AccAddress(sender.Bytes())}
}

// AsTransaction extracts the transaction as an `coretypes.Transaction`.
func (etr *WrappedEthereumTransaction) AsTransaction() *coretypes.Transaction {
	tx := new(coretypes.Transaction)
	if err := tx.UnmarshalBinary(etr.Data); err != nil {
		return nil
	}
	return tx
}

// GetSignBytes returns the bytes to sign over for the transaction.
func (etr *WrappedEthereumTransaction) GetSignBytes() ([]byte, error) {
	tx := etr.AsTransaction()
	return coretypes.LatestSignerForChainID(tx.ChainId()).
		Hash(tx).Bytes(), nil
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (etr *WrappedEthereumTransaction) GetSender() (common.Address, error) {
	tx := etr.AsTransaction()
	signer := coretypes.LatestSignerForChainID(tx.ChainId())
	return signer.Sender(tx)
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (etr *WrappedEthereumTransaction) GetPubKey() ([]byte, error) {
	tx := etr.AsTransaction()
	signer := coretypes.LatestSignerForChainID(tx.ChainId())
	return signer.PubKey(tx)
}

// GetSender extracts the sender address from the signature values using the latest signer for the given chainID.
func (etr *WrappedEthereumTransaction) GetSignature() ([]byte, error) {
	tx := etr.AsTransaction()
	signer := coretypes.LatestSignerForChainID(tx.ChainId())
	return signer.Signature(tx)
}

// GetGas returns the gas limit of the transaction.
func (etr *WrappedEthereumTransaction) GetGas() uint64 {
	var tx *coretypes.Transaction
	if tx = etr.AsTransaction(); tx == nil {
		return 0
	}
	return tx.Gas()
}

// GetGasPrice returns the gas price of the transaction.
func (etr *WrappedEthereumTransaction) ValidateBasic() error {
	// Ensure the transaction is signed properly
	tx := etr.AsTransaction()
	if tx == nil {
		return errors.New("transaction data is invalid")
	}

	// Ensure the transaction does not have a negative value.
	if tx.Value().Sign() < 0 {
		return txpool.ErrNegativeValue
	}

	// Sanity check for extremely large numbers.
	if tx.GasFeeCap().BitLen() > 256 { //nolint:gomnd // 256 bits.
		return core.ErrFeeCapVeryHigh
	}

	// Sanity check for extremely large numbers.
	if tx.GasTipCap().BitLen() > 256 { //nolint:gomnd // 256 bits.
		return core.ErrTipVeryHigh
	}

	// Ensure gasFeeCap is greater than or equal to gasTipCap.
	if tx.GasFeeCapIntCmp(tx.GasTipCap()) < 0 {
		return core.ErrTipAboveFeeCap
	}

	return nil
}

func (etr *WrappedEthereumTransaction) FeeGranter() sdk.AccAddress {
	return sdk.AccAddress{}
}

func (etr *WrappedEthereumTransaction) GetFee() sdk.Coins {
	return sdk.Coins{}
}

func (etr *WrappedEthereumTransaction) FeePayer() sdk.AccAddress {
	return sdk.AccAddress{}
}

func (etr *WrappedEthereumTransaction) GetMsgs() []sdk.Msg {
	return []sdk.Msg{etr}
}

// func (etr *WrappedEthereumTransaction) ToProto() client.TxBuilder {
// 	authtx.WrapTx(etr)
// 	return nil
// }

// GetAsEthTx is a helper function to get an EthTx from a sdk.Tx.
func GetAsEthTx(tx sdk.Tx) *coretypes.Transaction {
	if len(tx.GetMsgs()) == 0 {
		return nil
	}
	etr, ok := utils.GetAs[*WrappedEthereumTransaction](tx.GetMsgs()[0])
	if !ok {
		return nil
	}
	return etr.AsTransaction()
}

var _ client.TxEncodingConfig = (*EthTxEncodingConfig)(nil)

type EthTxEncodingConfig struct {
	signer coretypes.Signer
	cdc    codec.Codec
}

func NewEthEncodingConfig(chainID *big.Int) *EthTxEncodingConfig {
	return &EthTxEncodingConfig{
		signer: coretypes.LatestSignerForChainID(chainID),
	}
}

// Encoder
func (cfg EthTxEncodingConfig) TxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		// If its an ethereum transaction, return the data
		if ethTx, ok := tx.(*WrappedEthereumTransaction); ok {
			return ethTx.Data, nil
		}

		// else its a standard sdk.Tx, so just use the default txEncoder.
		return authtx.DefaultTxEncoder()(tx)
	}
}

func (cfg EthTxEncodingConfig) TxJSONEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		if ethTx, ok := tx.(*WrappedEthereumTransaction); ok {
			bz, err := ethTx.AsTransaction().MarshalJSON()
			if err == nil {
				return bz, nil
			}
		}
		return authtx.DefaultJSONTxEncoder(cfg.cdc)(tx)
	}
}

// Decoder
func (cfg EthTxEncodingConfig) TxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var ethTx coretypes.Transaction
		if err := ethTx.UnmarshalBinary(txBytes); err == nil {
			return &WrappedEthereumTransaction{Data: txBytes}, nil
		}

		return authtx.DefaultTxDecoder(cfg.cdc)(txBytes)
	}
}

func (cfg EthTxEncodingConfig) TxJSONDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var ethTx coretypes.Transaction
		if err := ethTx.UnmarshalJSON(txBytes); err == nil {
			return &WrappedEthereumTransaction{Data: txBytes}, nil
		}

		return authtx.DefaultJSONTxDecoder(cfg.cdc)(txBytes)
	}
}

// Signatures

func (g EthTxEncodingConfig) MarshalSignatureJSON(sigs []signing.SignatureV2) ([]byte, error) {
	descs := make([]*signing.SignatureDescriptor, len(sigs))

	for i, sig := range sigs {
		descData := signing.SignatureDataToProto(sig.Data)
		any, err := codectypes.NewAnyWithValue(sig.PubKey)
		if err != nil {
			return nil, err
		}

		descs[i] = &signing.SignatureDescriptor{
			PublicKey: any,
			Data:      descData,
			Sequence:  sig.Sequence,
		}
	}

	toJSON := &signing.SignatureDescriptors{Signatures: descs}

	return codec.ProtoMarshalJSON(toJSON, nil)
}

func (g EthTxEncodingConfig) UnmarshalSignatureJSON(bz []byte) ([]signing.SignatureV2, error) {
	var sigDescs signing.SignatureDescriptors
	err := g.cdc.UnmarshalJSON(bz, &sigDescs)
	if err != nil {
		return nil, err
	}

	sigs := make([]signing.SignatureV2, len(sigDescs.Signatures))
	for i, desc := range sigDescs.Signatures {
		pubKey, _ := desc.PublicKey.GetCachedValue().(cryptotypes.PubKey)

		data := signing.SignatureDataFromProto(desc.Data)

		sigs[i] = signing.SignatureV2{
			PubKey:   pubKey,
			Data:     data,
			Sequence: desc.Sequence,
		}
	}

	return sigs, nil
}

func (g EthTxEncodingConfig) NewTxBuilder() client.TxBuilder {
	return nil
}

func (g EthTxEncodingConfig) WrapTxBuilder(sdk.Tx) (client.TxBuilder, error) {
	return nil, nil
}

func (g EthTxEncodingConfig) SignModeHandler() *txsigning.HandlerMap {
	return nil
}
