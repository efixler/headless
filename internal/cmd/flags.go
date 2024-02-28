package cmd

import (
	"flag"

	"github.com/efixler/envflags"
)

type HeadlessBrowserSpec struct {
	Address *envflags.Value[string]
	Port    *envflags.Value[int]
}

func HeadlessBrowserFlags(flags *flag.FlagSet) *HeadlessBrowserSpec {
	hb := &HeadlessBrowserSpec{
		Address: envflags.NewString("BROWSER_ADDRESS", "127.0.0.1"),
		Port:    envflags.NewInt("BROWSER_PORT", 9222),
	}
	hb.addToFlagSet(flags)
	return hb
}

func (hb *HeadlessBrowserSpec) addToFlagSet(flags *flag.FlagSet) {
	hb.Address.AddTo(flags, "browser-address", "Address of the headless browser")
	hb.Port.AddTo(flags, "browser-port", "Port of the headless browser")
}
