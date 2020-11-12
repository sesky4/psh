package main

import (
	"fmt"
	"log"
	"os"
	"psh"
)

func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c := psh.LoadConfig()
	ifAutoComplete(c)
	var host string
	var args []string

	// 判断传入ssh的参数
	// 1. "-Xxxx" 的形式
	// 2. "-X xxx" 的形式
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg[0] == '-' {
			if len(arg) < 2 {
				fmt.Fprintf(os.Stderr, "unknown arg "+arg)
			} else if len(arg) == 2 {
				// -X xxx的形式
				args = append(args, arg, os.Args[i+1])
				i++
			} else if len(arg) > 2 {
				// -Xxxx的形式
				args = append(args, arg)
			}
		} else {
			host = arg
		}
	}

	item, ok := c.Match(host)
	if !ok {
		log.Fatalf("%s not found in ~/.ssh/pconfig", host)
	}
	proxySSH(item.Host, item.Port, item.Username, item.Passwords, args)
}
