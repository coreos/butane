package config

import (
	base_0_1 "github.com/ajeddeloh/fcct/base/v0_1"
	"github.com/ajeddeloh/fcct/distro/fcos_0_1"
)

type FcosConfig0_1 struct {
	Common          `yaml:",inline"`
	base_0_1.Config `yaml:",inline"`
	fcos_0_1.Fcos   `yaml:",inline"`
}

func dumbMerge(a, b interface{}) {
	// pass
}

func TranslateFcos0_1(input []byte, options TranslateOptions) ([]byte, error) {
	cfg := FcosConfig0_1{}

	if err := unmarshal(input, &cfg, options.Strict); err != nil {
		return nil, err
	}

	base, _ := cfg.Config.ToIgn3_0()
	distro, _ := cfg.Fcos.ToIgn3_0()
	// do a dumb merge, these should not conflict and if they do the user should fix them
	dumbMerge(&base, distro)

	// TODO validation
	return marshal(base, options.Pretty)
}
