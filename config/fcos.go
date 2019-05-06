package config

import (
	"errors"
	"fmt"
	"reflect"

	base_0_1 "github.com/ajeddeloh/fcct/base/v0_1"
	fcos_0_1 "github.com/ajeddeloh/fcct/distro/fcos/v0_1"

	"github.com/coreos/ignition/v2/config/v3_0"
	"github.com/coreos/ignition/v2/config/validate"
)

var (
	ErrInvalidConfig = errors.New("config generated was invalid")
)

type FcosConfig0_1 struct {
	Common          `yaml:",inline"`
	base_0_1.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func TranslateFcos0_1(input []byte, options TranslateOptions) ([]byte, error) {
	cfg := FcosConfig0_1{}

	if err := unmarshal(input, &cfg, options.Strict); err != nil {
		return nil, err
	}

	base, err := cfg.Config.ToIgn3_0()
	if err != nil {
		return nil, err
	}

	distro, err := cfg.Fcos.ToIgn3_0()
	if err != nil {
		return nil, err
	}

	final := v3_0.Merge(distro, base)
	r := validate.ValidateWithoutSource(reflect.ValueOf(final))
	fmt.Println(r.String())
	if r.IsFatal() {
		return nil, ErrInvalidConfig
	}

	// TODO validation
	return marshal(final, options.Pretty)
}
