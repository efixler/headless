package cmd

import (
	"flag"

	"github.com/efixler/envflags"
)

type HeadlessBrowserSpec struct {
	Address *envflags.Value[string]
	Port    *envflags.Value[int]
}

func HeadlessBrowserFlags(baseEnv string, flags *flag.FlagSet) *HeadlessBrowserSpec {
	hb := &HeadlessBrowserSpec{
		Address: envflags.NewString(baseEnv+"BROWSER_ADDRESS", "127.0.0.1"),
		Port:    envflags.NewInt(baseEnv+"BROWSER_PORT", 9222),
	}
	hb.addToFlagSet(flags)
	return hb
}

func (hb *HeadlessBrowserSpec) addToFlagSet(flags *flag.FlagSet) {
	flags.Var(hb.Address, "browser-address", "Address of the headless browser")
	flags.Var(hb.Port, "browser-port", "Port of the headless browser")
}
