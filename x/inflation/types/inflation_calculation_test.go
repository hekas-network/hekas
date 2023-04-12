package types

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InflationTestSuite struct {
	suite.Suite
}

func TestInflationSuite(t *testing.T) {
	suite.Run(t, new(InflationTestSuite))
}

func (suite *InflationTestSuite) TestCalculateEpochMintProvision() {
	bondingParams := DefaultParams()
	bondingParams.ExponentialCalculation.MaxVariance = sdk.NewDecWithPrec(40, 2)
	epochsPerPeriod := int64(365)

	testCases := []struct {
		name              string
		params            Params
		period            uint64
		bondedRatio       sdk.Dec
		expEpochProvision sdk.Dec
		expPass           bool
	}{
		{
			"pass - initial perid",
			DefaultParams(),
			uint64(0),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("17643835616438356164384.000000000000000000"),
			true,
		},
		{
			"pass - period 1",
			DefaultParams(),
			uint64(1),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("9972602739726027397260.000000000000000000"),
			true,
		},
		{
			"pass - period 2",
			DefaultParams(),
			uint64(2),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("6136986301369863013699.000000000000000000"),
			true,
		},
		{
			"pass - period 3",
			DefaultParams(),
			uint64(3),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("4219178082191780821918.000000000000000000"),
			true,
		},
		{
			"pass - period 20",
			DefaultParams(),
			uint64(20),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("2301384494729238009863.000000000000000000"),
			true,
		},
		{
			"pass - period 21",
			DefaultParams(),
			uint64(21),
			sdk.OneDec(),
			sdk.MustNewDecFromStr("2301377178871468320000.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - initial period",
			bondingParams,
			uint64(0),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("24701369863013698630137.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 1",
			bondingParams,
			uint64(1),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("13961643835616438356164.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 2",
			bondingParams,
			uint64(2),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("8591780821917808219178.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 3",
			bondingParams,
			uint64(3),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("5906849315068493150685.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 20",
			bondingParams,
			uint64(20),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("3221938292620933213808.000000000000000000"),
			true,
		},
		{
			"pass - 0 percent bonding - period 21",
			bondingParams,
			uint64(21),
			sdk.ZeroDec(),
			sdk.MustNewDecFromStr("3221928050420055648000.000000000000000000"),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			epochMintProvisions := CalculateEpochMintProvision(
				tc.params,
				tc.period,
				epochsPerPeriod,
				tc.bondedRatio,
			)

			suite.Require().Equal(tc.expEpochProvision, epochMintProvisions)
		})
	}
}
