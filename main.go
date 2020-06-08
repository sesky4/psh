package main

import (
	"log"
	"os"
	"psh/conf"
)

func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c := conf.LoadConfig()
	ifAutoComplete(c)
	host := os.Args[1]

	item, ok := c.Match(host)
	if !ok {
		log.Fatalf("%s not found in ~/.ssh/pconfig", host)
	}
	proxySSH(item.Host, item.Port, item.Username, item.Passwords)
}
