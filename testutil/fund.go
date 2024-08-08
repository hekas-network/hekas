// Copyright 2022 Hekas Foundation
// This file is part of the Hekas Network packages.
//
// Hekas is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Hekas packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Hekas packages. If not, see https://github.com/hekas-network/hekas/blob/main/LICENSE

package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/hekas-network/hekas/utils"
	inflationtypes "github.com/hekas-network/hekas/x/inflation/types"
)

// FundAccount is a utility function that funds an account by minting and
// sending the coins to the address.
func FundAccount(ctx sdk.Context, bankKeeper bankkeeper.Keeper, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, inflationtypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, inflationtypes.ModuleName, addr, amounts)
}

// FundAccountWithBaseDenom is a utility function that uses the FundAccount function
// to fund an account with the default Hekas denomination.
func FundAccountWithBaseDenom(ctx sdk.Context, bankKeeper bankkeeper.Keeper, addr sdk.AccAddress, amount int64) error {
	coins := sdk.NewCoins(
		sdk.NewCoin(utils.BaseDenom, sdk.NewInt(amount)),
	)
	return FundAccount(ctx, bankKeeper, addr, coins)
}

// FundModuleAccount is a utility function that funds a module account by
// minting and sending the coins to the address.
func FundModuleAccount(ctx sdk.Context, bankKeeper bankkeeper.Keeper, recipientMod string, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, inflationtypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToModule(ctx, inflationtypes.ModuleName, recipientMod, amounts)
}
