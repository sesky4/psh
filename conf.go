package psh

import (
	"fmt"
	"github.com/gobwas/glob"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type Conf struct {
	Named map[string]LoginInfo
	Globs []LoginInfo
}

func (c *Conf) Register(info LoginInfo) {
	if IsGlob(info.Pattern) {
		c.Globs = append(c.Globs, info)
	} else {
		c.Named[info.Pattern] = info
	}
}

// 调用Register的Order会影响Match顺序
func (c *Conf) Match(pattern string) (LoginInfo, bool) {
	if info, ok := c.Named[pattern]; ok {
		return info, ok
	} else {
		for _, info := range c.Globs {
			if glob.MustCompile(info.Pattern).Match(pattern) {
				return info, true
			}
		}
		return LoginInfo{}, false
	}
}

func IsGlob(pattern string) bool {
	return strings.IndexAny(pattern, "*[]{}?!") >= 0
}

type LoginInfo struct {
	Pattern   string
	Username  string
	Passwords []string
	Host      string
	Port      int
}

func LoadConfig() Conf {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	cPath := filepath.Join(u.HomeDir, ".ssh", "pconfig")
	c, err := ioutil.ReadFile(cPath)
	if err != nil {
		panic(err)
	}

	config := Conf{
		Named: map[string]LoginInfo{},
	}
	lines := strings.Split(string(c), "\n")

	// format:
	// ${name}\t${user}@${host}:${port}\t${pwd}
	// eg:
	// 91	root@192.168.199.91	123456
	for lineNo, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// 注释
		if line[0] == '#' {
			continue
		}

		info, err := parseConfLine(line)
		if err != nil {
			log.Println(line)
			log.Fatalf("~/.ssh/pconfig:%d   %s", lineNo, err)
		}

		config.Register(info)
	}
	return config
}

func parseConfLine(line string) (info LoginInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	bs := strings.Split(line, "\t")
	if len(bs) != 3 {
		err = fmt.Errorf("line split by '\\t' is not 3 segements")
		return
	}

	var username, host string
	var port int
	var pwds []string

	pattern := bs[0]

	s2 := strings.Split(bs[1], "@")
	username = s2[0]
	s22 := strings.Split(s2[1], ":")
	host = s22[0]
	if len(s22) == 1 {
		port = 22
	} else if len(s22) == 2 {
		port, err = strconv.Atoi(s22[1])
	}

	pwds = strings.Split(bs[2], " ")

	info.Pattern = pattern
	info.Host = host
	info.Port = port
	info.Username = username
	info.Passwords = pwds
	return
}
