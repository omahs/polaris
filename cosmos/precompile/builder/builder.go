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

package auth

import (
	"context"
	"math/big"

	generated "pkg.berachain.dev/polaris/contracts/bindings/cosmos/precompile"
	"pkg.berachain.dev/polaris/cosmos/precompile"
	"pkg.berachain.dev/polaris/eth/common"
	ethprecompile "pkg.berachain.dev/polaris/eth/core/precompile"
	"pkg.berachain.dev/polaris/eth/params"
)

// Contract is the precompile contract for the auth module.
type Contract struct {
	precompile.BaseContract
}

// NewPrecompileContract returns a new instance of the auth module precompile contract.
func NewPrecompileContract( /*BuilderKeeper*/ ) ethprecompile.StatefulImpl {
	return &Contract{
		BaseContract: precompile.NewBaseContract(
			generated.BuilderModuleMetaData.ABI,
			// Replace 123 with the commented out below.
			common.Address{123},
			// cosmlib.AccAddressToEthAddress(
			// 	authtypes.NewModuleAddress(buildertype.ModuleName),
			// ),
		),
	}
}

// PrecompileMethods implements StatefulImpl.
func (c *Contract) PrecompileMethods() ethprecompile.Methods {
	return ethprecompile.Methods{
		{
			AbiSig:      "auctionBid(string,Coin,bytes)",
			Execute:     c.AuctionBid,
			RequiredGas: params.IdentityBaseGas,
		},
	}
}

func (c *Contract) AuctionBid(
	ctx context.Context,
	_ ethprecompile.EVM,
	caller common.Address,
	value *big.Int,
	readonly bool,
	args ...any,
) ([]any, error) {
	// 1. deserialize bytes
	// 2. build MsgAuctionBid
	// 3. C.msgServer.HandleMsgAuctionBid(.....)
	return []any{}, nil
}
