// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	stakingprecompile "github.com/evmos/evmos/v18/precompiles/staking"
	"github.com/evmos/evmos/v18/x/evm/types"
)

func (suite *KeeperTestSuite) TestIsAvailablePrecompile() {
	testcases := []struct {
		name         string
		address      common.Address
		expAvailable bool
	}{
		{
			name:         "pass - available precompile",
			address:      common.HexToAddress(stakingprecompile.PrecompileAddress),
			expAvailable: true,
		},
		{
			name:         "fail - unavailable precompile",
			address:      common.HexToAddress("0x0000000000000000000000000000000000099999"),
			expAvailable: false,
		},
	}

	for _, tc := range testcases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			available := suite.app.EvmKeeper.IsAvailablePrecompile(tc.address)
			suite.Require().Equal(tc.expAvailable, available)
		})
	}
}

// Check interface is correctly implemented
var (
	_ vm.PrecompiledContract = DummyPrecompile{}
)

// DummyPrecompile is a dummy precompile implementation for testing purposes.
type DummyPrecompile struct {
	vm.PrecompiledContract

	address string
}

func (d DummyPrecompile) Address() common.Address {
	return common.HexToAddress(d.address)
}

var (
	// dummyPrecompile holds an unused precompile address to check adding EVM extensions.
	dummyPrecompile = DummyPrecompile{address: "0x0000000000000000000000000000000000010000"}
	// otherPrecompile holds another unused precompile address to check adding multiple extensions at once.
	otherPrecompile = DummyPrecompile{address: "0x0000000000000000000000000000000000010001"}
)

func (suite *KeeperTestSuite) TestAddDynamicEVMExtensions() {
	testcases := []struct {
		name           string
		malleate       func() []vm.PrecompiledContract
		expPass        bool
		errContains    string
		expPrecompiles []string
	}{
		{
			name: "fail - add multiple precompiles with duplicates",
			malleate: func() []vm.PrecompiledContract {
				return []vm.PrecompiledContract{dummyPrecompile, dummyPrecompile}
			},
			errContains: "duplicate precompile",
		},
		{
			name: "fail - precompile already in active precompiles",
			malleate: func() []vm.PrecompiledContract {
				// We add the dummy precompile to the active dynamic precompiles to trigger the error.
				err := suite.app.EvmKeeper.AddDynamicPrecompiles(suite.ctx, dummyPrecompile)
				suite.Require().NoError(err, "expected no error adding extensions")

				return []vm.PrecompiledContract{dummyPrecompile}
			},
			errContains: "precompile already registered",
		},
		{
			name: "pass - add precompile",
			malleate: func() []vm.PrecompiledContract {
				return []vm.PrecompiledContract{dummyPrecompile}
			},
			expPass:        true,
			expPrecompiles: append(types.AvailableEVMExtensions, dummyPrecompile.Address().String()),
		},
		{
			name: "pass - add multiple precompiles",
			malleate: func() []vm.PrecompiledContract {
				return []vm.PrecompiledContract{dummyPrecompile, otherPrecompile}
			},
			expPass:        true,
			expPrecompiles: append(types.AvailableEVMExtensions, dummyPrecompile.Address().String(), otherPrecompile.Address().String()),
		},
	}

	for _, tc := range testcases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			var extensions []vm.PrecompiledContract
			suite.Require().NotNil(tc.malleate, "malleate must be defined")
			extensions = tc.malleate()

			err := suite.app.EvmKeeper.AddDynamicPrecompiles(suite.ctx, extensions...)
			if tc.expPass {
				suite.Require().NoError(err, "expected no error adding extensions")

				activeStaticPrecompiles := suite.app.EvmKeeper.GetParams(suite.ctx).ActiveStaticPrecompiles
				activeDynamicPrecompiles := suite.app.EvmKeeper.GetParams(suite.ctx).ActiveDynamicPrecompiles
				combinedPrecompiles := append(activeStaticPrecompiles, activeDynamicPrecompiles...) //nolint:gocritic // use of append is fine here
				suite.Require().Equal(tc.expPrecompiles, combinedPrecompiles, "expected different active precompiles")

				for _, expPrecompile := range tc.expPrecompiles {
					suite.Require().Contains(combinedPrecompiles, expPrecompile, "expected available precompiles to contain: %s", expPrecompile)
				}
			} else {
				suite.Require().Error(err, "expected error adding extensions")
				suite.Require().ErrorContains(err, tc.errContains, "expected different error")
			}
		})
	}
}
