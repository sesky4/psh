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
	"syscall"
)

var (
	pwdRE      = regexp.MustCompile(`.*password:`)
	certRE     = regexp.MustCompile(`The authenticity of host .* can't be established.*`)
	wrongPwdRE = regexp.MustCompile(`Permission denied, please try again.`)
	loginRE    = regexp.MustCompile(`.*Welcome to .*`)
)

func proxySCPPull(rUser, rHost, rPath string, rPort int, lPath string, passwords []string) {
	// PULL: scp -C -Pport root@192.168.0.1:/home/file ./file.py
	args := []string{"-C", fmt.Sprintf("-P%d", rPort), fmt.Sprintf("%s@%s:%s", rUser, rHost, rPath), lPath}
	cmd := exec.Command("scp", args...)

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

	// login
	{
		buf := make([]byte, 30*1024)
		for {
			nr, err := ptmx.Read(buf)
			poe(err)
			os.Stdout.Write(buf[:nr])

			rs := string(buf[:nr])
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
	go io.Copy(ptmx, os.Stdin)
	io.Copy(os.Stdout, ptmx)
}

func proxySCPPush(rUser, rHost, rPath string, rPort int, lPath string, passwords []string) {
	// PUSH: scp -C -Pport ./file.py root@192.168.0.1:/home/file
	args := []string{"-C", fmt.Sprintf("-P%d", rPort), lPath, fmt.Sprintf("%s@%s:%s", rUser, rHost, rPath)}
	cmd := exec.Command("scp", args...)
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

	// login
	{
		buf := make([]byte, 30*1024)
		for {
			nr, err := ptmx.Read(buf)
			poe(err)
			os.Stdout.Write(buf[:nr])

			rs := string(buf[:nr])
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
	go io.Copy(ptmx, os.Stdin)
	io.Copy(os.Stdout, ptmx)
}
