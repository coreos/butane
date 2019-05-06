package v0_1

import (
	"github.com/coreos/ignition/v2/config/v3_0/types"
)

func (c Config) ToIgn3_0() (types.Config, error) {
	return types.Config{
		Ignition: types.Ignition{
			Version: "3.0.0",
		},
	}, nil
}
