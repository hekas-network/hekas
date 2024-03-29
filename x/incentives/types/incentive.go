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

package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	hekastypes "github.com/hekas-network/hekas/types"
)

// NewIncentive returns an instance of Incentive
func NewIncentive(
	contract common.Address,
	allocations sdk.DecCoins,
	epochs uint32,
) Incentive {
	return Incentive{
		Contract:    contract.String(),
		Allocations: allocations,
		Epochs:      epochs,
		TotalGas:    0,
	}
}

// Validate performs a stateless validation of a Incentive
func (i Incentive) Validate() error {
	if err := hekastypes.ValidateAddress(i.Contract); err != nil {
		return err
	}

	if i.Allocations.IsZero() {
		return fmt.Errorf("allocations cannot be empty: %s", i.Allocations)
	}

	for _, al := range i.Allocations {
		if err := sdk.ValidateDenom(al.Denom); err != nil {
			return err
		}
		if err := validateAmount(al.Amount); err != nil {
			return err
		}
	}

	if i.Epochs == 0 {
		return fmt.Errorf("epoch cannot be 0")
	}
	return nil
}

// IsActive returns true if the Incentive has remaining Epochs
func (i Incentive) IsActive() bool {
	return i.Epochs > 0
}
