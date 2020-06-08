package main

import (
	"flag"
	"github.com/posener/complete"
	"os"
)

// bash/zsh auto-completion
func ifAutoComplete(c Conf) {
	cmd := complete.Command{
		Flags: complete.Flags{},
	}
	for k, _ := range c.Named {
		cmd.Flags[k] = complete.PredictAnything
	}
	cmp := complete.New(
		"psh",
		cmd,
	)

	// AddFlags adds the completion flags to the program flags,
	// in case of using non-default flag set, it is possible to pass
	// it as an argument.
	// it is possible to set custom flags name
	// so when one will type 'self -h', he will see '-complete' to install the
	// completion and -uncomplete to uninstall it.
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)

	flag.Parse()
	if cmp.Complete() {
		os.Exit(0)
	}
}
