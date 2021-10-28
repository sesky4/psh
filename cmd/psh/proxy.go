package main

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
)

var (
	pwdRE      = regexp.MustCompile(`.*password:`)
	certRE     = regexp.MustCompile(`The authenticity of host .* can't be established.*`)
	wrongPwdRE = regexp.MustCompile(`Permission denied, please try again.`)
	loginRE    = regexp.MustCompile(`.*Last login.*`)
)

func proxySSH(host string, port int, user string, passwords, args []string) {
	cmdArgs := []string{"-p", strconv.Itoa(port), fmt.Sprintf("%s@%s", user, host)}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("ssh", cmdArgs...)

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	poe(err)
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	go io.Copy(ptmx, os.Stdin)
	// login
	{
		buf := make([]byte, 30*1024)
		for {
			nr, err := ptmx.Read(buf)
			poe(err)

			bs := buf[:nr]
			os.Stdout.Write(bs)

			rs := string(bs)
			switch {
			case pwdRE.MatchString(rs):
				ptmx.Write([]byte(passwords[0] + "\n"))
			case certRE.MatchString(rs):
				ptmx.Write([]byte("yes\n"))
			case wrongPwdRE.MatchString(rs):
				passwords = passwords[1:]
			case loginRE.MatchString(rs):
				goto loginSuccess
			}
		}
	}

loginSuccess:
	io.Copy(os.Stdout, ptmx)
}
