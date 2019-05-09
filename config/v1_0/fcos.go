package v1_0

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ajeddeloh/fcct/config/common"
	base_0_1 "github.com/ajeddeloh/fcct/base/v0_1"
	fcos_0_1 "github.com/ajeddeloh/fcct/distro/fcos/v0_1"

	"github.com/coreos/ignition/v2/config/v3_0"
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/ignition/v2/config/validate"
)

var (
	ErrInvalidConfig = errors.New("config generated was invalid")
)

type Config struct {
	common.Common   `yaml:",inline"`
	base_0_1.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func (c Config) Translate() (types.Config, error) {
	base, err := c.Config.ToIgn3_0()
	if err != nil {
		return types.Config{}, err
	}

	distro, err := c.Fcos.ToIgn3_0()
	if err != nil {
		return types.Config{}, err
	}

	return v3_0.Merge(distro, base), nil
}

func TranslateBytes(input []byte, options common.TranslateOptions) ([]byte, error) {
	cfg := Config{}

	if err := common.Unmarshal(input, &cfg, options.Strict); err != nil {
		return nil, err
	}
	r := validate.ValidateWithoutSource(reflect.ValueOf(cfg))
	if r.IsFatal() {
		fmt.Println(r.String())
		return nil, ErrInvalidConfig
	}

	final, err := cfg.Translate()
	if err != nil {
		return nil, err
	}

	r.Merge(validate.ValidateWithoutSource(reflect.ValueOf(final)))
	fmt.Println(r.String())

	if r.IsFatal() {
		return nil, ErrInvalidConfig
	}

	return common.Marshal(final, options.Pretty)
}
