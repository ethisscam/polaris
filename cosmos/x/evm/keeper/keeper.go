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

package keeper

import (
	"context"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/berachain/polaris/cosmos/config"
	"github.com/berachain/polaris/cosmos/x/evm/plugins/state"
	"github.com/berachain/polaris/cosmos/x/evm/types"
	"github.com/berachain/polaris/eth/core"
	ethprecompile "github.com/berachain/polaris/eth/core/precompile"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type WrappedBlockchain interface {
	PreparePlugins(context.Context)
	Config() *params.ChainConfig
	WriteGenesisState(context.Context, *core.Genesis) error
	InsertBlockAndSetHead(context.Context, *ethtypes.Block) error
	GetBlockByNumber(uint64) *ethtypes.Block
}

type Keeper struct {
	// host represents the host chain
	*Host

	// provider is the struct that houses the Polaris EVM.
	wrappedChain WrappedBlockchain
}

// NewKeeper creates new instances of the polaris Keeper.
func NewKeeper(
	ak state.AccountKeeper,
	storeKey storetypes.StoreKey,
	pcs func() *ethprecompile.Injector,
	qc func() func(height int64, prove bool) (sdk.Context, error),
	polarisCfg *config.Config,
) *Keeper {
	host := NewHost(
		*polarisCfg,
		storeKey,
		ak,
		pcs,
		qc,
	)
	return &Keeper{
		Host: host,
	}
}

func (k *Keeper) Setup(wrappedChain WrappedBlockchain) error {
	k.wrappedChain = wrappedChain
	return k.SetupPrecompiles()
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(types.ModuleName)
}
