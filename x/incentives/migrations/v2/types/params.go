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
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/hekas-network/hekas/x/incentives/types"

	epochstypes "github.com/hekas-network/hekas/x/epochs/types"
)

var _ types.LegacyParams = &V2Params{}

// Parameter store key
var (
	ParamsKey                     = []byte("Params")
	ParamStoreKeyEnableIncentives = []byte("EnableIncentives")
	ParamStoreKeyAllocationLimit  = []byte("AllocationLimit")
	ParamStoreKeyEpochIdentifier  = []byte("EpochIdentifier")
	ParamStoreKeyRewardScaler     = []byte("RewardScaler")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&V2Params{})
}

var (
	DefaultEnableIncentives          = true
	DefaultAllocationLimit           = types.DefaultAllocationLimit
	DefaultIncentivesEpochIdentifier = epochstypes.WeekEpochID
	DefaultRewardScalar              = types.DefaultRewardScalar
)

// NewParams creates a new Params object
func NewParams(
	enableIncentives bool,
	allocationLimit sdk.Dec,
	epochIdentifier string,
	rewardScaler sdk.Dec,
) V2Params {
	return V2Params{
		EnableIncentives:          enableIncentives,
		AllocationLimit:           allocationLimit,
		IncentivesEpochIdentifier: epochIdentifier,
		RewardScaler:              rewardScaler,
	}
}

func DefaultParams() V2Params {
	return V2Params{
		EnableIncentives:          DefaultEnableIncentives,
		AllocationLimit:           DefaultAllocationLimit,
		IncentivesEpochIdentifier: DefaultIncentivesEpochIdentifier,
		RewardScaler:              DefaultRewardScalar,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *V2Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableIncentives, &p.EnableIncentives, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyAllocationLimit, &p.AllocationLimit, validatePercentage),
		paramtypes.NewParamSetPair(ParamStoreKeyEpochIdentifier, &p.IncentivesEpochIdentifier, epochstypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(ParamStoreKeyRewardScaler, &p.RewardScaler, validateUncappedPercentage),
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePercentage(i interface{}) error {
	dec, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if dec.IsNil() {
		return errors.New("allocation limit cannot be nil")
	}
	if dec.IsNegative() {
		return fmt.Errorf("allocation limit must be positive: %s", dec)
	}
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("allocation limit must <= 100: %s", dec)
	}

	return nil
}

func validateUncappedPercentage(i interface{}) error {
	dec, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if dec.IsNil() {
		return errors.New("allocation limit cannot be nil")
	}
	if dec.IsNegative() {
		return fmt.Errorf("allocation limit must be positive: %s", dec)
	}

	return nil
}

func (p V2Params) Validate() error {
	if err := validateBool(p.EnableIncentives); err != nil {
		return err
	}

	if err := validatePercentage(p.AllocationLimit); err != nil {
		return err
	}

	if err := validateUncappedPercentage(p.RewardScaler); err != nil {
		return err
	}

	return epochstypes.ValidateEpochIdentifierString(p.IncentivesEpochIdentifier)
}