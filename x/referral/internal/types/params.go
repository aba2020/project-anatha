package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName
)

var (
	KeyNameReferralPercentage	= []byte("ReferralPercentage")

	DefaultReferralPercentage = sdk.NewDecWithPrec(1, 1)
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	ReferralPercentage		sdk.Dec 	`json:"referral_percentage" yaml:"referral_percentage"`
}

func NewParams(referralPercentage sdk.Dec) Params {
	return Params{
		ReferralPercentage: referralPercentage,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Referral Percentage:     %s`,
		p.ReferralPercentage,
	)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyNameReferralPercentage, &p.ReferralPercentage, validateReferralPercentage),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultReferralPercentage,
	)
}

func (p Params) Validate() error {
	if err := validateReferralPercentage(p.ReferralPercentage); err != nil {
		return err
	}

	return nil
}

func validateReferralPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("referral percentage must be less than 100%: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("referral percentage must be positive: %s", v)
	}

	return nil
}

