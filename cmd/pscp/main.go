package main

import (
	"flag"
	"fmt"
	"github.com/posener/complete"
	"os"
	"psh"
	"strings"
)

func main() {
	c := psh.LoadConfig()
	args := make([]string, len(os.Args))
	copy(args, os.Args)

	var mode string
	var rUser, rHost, rPath, lPath string
	var rPort int
	var passwords []string

	// replace hostname while keep other options in place
	for _, arg := range args[1:] {
		if strings.Index(arg, ":") > 0 {
			if mode == "" {
				mode = "pull"
			}
			ss := strings.Split(arg, ":")
			host := ss[0]
			fpath := ss[1]
			info, ok := c.Match(host)
			if !ok {
				fmt.Printf("%s not found in ~/.ssh/pconfig\n", host)
				os.Exit(1)
			}

			rUser = info.Username
			rHost = info.Host
			rPath = fpath
			rPort = info.Port
			passwords = info.Passwords
		} else {
			if mode == "" {
				mode = "push"
			}
			lPath = arg
		}
	}

	if mode == "push" {
		proxySCPPush(rUser, rHost, rPath, rPort, lPath, passwords)
	} else {
		proxySCPPull(rUser, rHost, rPath, rPort, lPath, passwords)
	}
}

func poe(err error) {
	if err != nil {
		panic(err)
	}
}

// bash/zsh auto-completion
func ifAutoComplete(c psh.Conf) {
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
