package main

import (
	"flag"
	"fmt"
	"github.com/posener/complete"
	"os"
	main2 "psh"
)

func main() {
	_ := main2.LoadConfig()
	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])
	fmt.Println(os.Args[3])
}

// bash/zsh auto-completion
func ifAutoComplete(c main2.Conf) {
	cmd := complete.Command{
		Flags: complete.Flags{},
	}
	for k, _ := range c.Named {
		cmd.Flags[k] = complete.PredictAnything
	}
	cmp := complete.New(
		"pscp",
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
