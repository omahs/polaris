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

package historical

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testutil "pkg.berachain.dev/polaris/cosmos/testing/utils"
	"pkg.berachain.dev/polaris/eth/core/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Historical Plugin", func() {
	var (
		p   *plugin
		ctx sdk.Context
	)

	BeforeEach(func() {
		ctx = testutil.NewContext()
		p = &plugin{
			ctx:              ctx,
			bp:               mock.NewBlockPluginMock(),
			storeKey:         storetypes.NewKVStoreKey("evm"),
			offchainStoreKey: storetypes.NewKVStoreKey("offchain-evm"),
		}
	})

	Context("After Genesis", func() {
		When("BlockByNumber is called on block 0", func() {
			It("should return the header without error", func() {
				block, err := p.GetBlockByNumber(0)
				Expect(err).ToNot(HaveOccurred())
				header := block.Header()
				Expect(header).ToNot(BeNil()) // mock header
			})
		})
	})

	// It("should get the header at current height", func() {
	// 	header, err := p.GetHeaderByNumber(ctx.BlockHeight())
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(header.TxHash).To(Equal(common.BytesToHash(ctx.BlockHeader().DataHash)))
	// })

	// It("should return empty header for non-existent height", func() {
	// 	header, err := p.GetHeaderByNumber(100000)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(*header).To(Equal(types.Header{}))
	// })
})
