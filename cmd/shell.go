package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"
)

type SysInfo struct {
	Username, Hostname string
}

func GetUsername() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func GetSysinfo() (SysInfo, error) {
	si := SysInfo{}

	hostname, err := os.Hostname()
	if err != nil {
		return si, err
	}
	si.Hostname = hostname

	username, err := GetUsername()
	if err != nil {
		return si, err
	}
	si.Username = username

	return si, nil
}

// Type out a string
func SlowPrint(param string) {
	// Last character type. Set to -1 by default (matches no character)
	lastChar := -1
	for _, i := range param {
		// minMS describes the minimum amount of time to type the character
		// rangeMS describes the maximum amount of time added to this minimum value
		minMS, rangeMS := 17, 80
		if lastChar == ' ' {
			// If last character was a space, refocus on the keys
			minMS, rangeMS = 100, 75
		} else if lastChar != '-' && i == '-' {
			// Same with the hyphen
			minMS, rangeMS = 75, 50
		}

		// Calculate and wait for delay
		calcMS := rand.Float32()*float32(rangeMS) + float32(minMS)
		time.Sleep(time.Millisecond * time.Duration(calcMS))

		// Print the character
		fmt.Print(string(i))
	}
}

// Restore the screen to its prior state
func Restore() {
	fmt.Print("\x1b[?1049l")
}

// Main application function
func shlmain(opts CommandOptions) {
	inf, err := GetSysinfo()
	if err != nil {
		log.Fatalln("could not get system info:", err)
	}

	// *** Initialization ***
	// Enable alt buffer
	fmt.Print("\x1b[?1049h\x1b[2J\x1b[H")

	// catch SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		// upon SIGINT, restore the screen and exit
		Restore()
		os.Exit(0)
	}()

	// *** SHELL PART ***
	// print prompt
	fmt.Printf("[%s@%s ~]$ ", inf.Username, inf.Hostname)

	// reaction time
	time.Sleep(time.Second * 1)

	// type out the command
	rootProg := ""
	separator := " "
	baseCmd := "rm -rvf --no-preserve-root /"

	switch {
	case opts.UseSudo:
		rootProg = "sudo"
	case opts.UseDoas:
		rootProg = "doas"
	case opts.UseRun0:
		rootProg = "run0"
	case opts.UsePkExec:
		rootProg = "pkexec"
	default:
		rootProg = "sudo"
	}
	SlowPrint(rootProg + separator + baseCmd)

	// wait a lil
	time.Sleep(time.Millisecond * 400)

	// then hit enter
	fmt.Println("")

	// *** simulate auth ***
	Auth(rootProg, inf)

	// finally, execute the remove command
	RecurseOver("/")

	// After the command, prompt and wait for C-c
	fmt.Printf("[%s@%s ~]$ ", inf.Username, inf.Hostname)
	for {
		time.Sleep(time.Second)
	}
}
