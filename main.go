package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"os/user"
	"path"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
)

type SysInfo struct {
	Username, Hostname string
}

type CommandOptions struct {
	UseSudo, UseDoas, UseRun0, UsePkExec bool
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

// Executes the remove command
func RecurseOver(dir string) {
	// Stat the path so we can see what it is.
	time.Sleep(time.Millisecond)
	stat, err := os.Lstat(dir)

	switch {
	// If an error happens, display it in a perror-style way
	case err != nil:
		fmt.Printf("rm: cannot remove %s: %s\n", dir, err)

	// If it's a directory, recurse over its children and delete them as well.
	case stat.IsDir():
		listing, err := os.ReadDir(dir)
		slices.SortFunc(listing, func(a, b os.DirEntry) int {
			return strings.Compare(a.Name(), b.Name())
		})
		if err != nil {
			fmt.Printf("rm: cannot remove %s: %s\n", dir, err)
			return
		}
		for _, i := range listing {
			RecurseOver(path.Join(dir, i.Name()))
		}
		fmt.Printf("removed directory '%s'\n", dir)

	// Symbolic links are treated specially.
	case stat.Mode()&os.ModeSymlink != 0:
		fmt.Printf("removed symbolic link '%s'\n", dir)

	// Otherwise assume it's a regular file.
	default:
		fmt.Printf("removed '%s'\n", dir)
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
	path, action := "", ""
	baseCmd := "rm -rvf --no-preserve-root /"

	switch {
	case opts.UseSudo:
		rootProg = "sudo"
	case opts.UseDoas:
		rootProg = "doas"
	case opts.UseRun0:
		rootProg = "run0"
		path = "org.freedesktop.systemd1.manage-units"
		action = "start transient unit 'run-p989-i990.service'."
	case opts.UsePkExec:
		rootProg = "pkexec"
		path = "org.freedesktop.policykit.exec"
		action = "run '/bin/rm' as the super user"
	default:
		rootProg = "sudo"
	}
	SlowPrint(rootProg + separator + baseCmd)

	// wait a lil
	time.Sleep(time.Millisecond * 400)

	// then hit enter
	fmt.Println("")

	// *** SUDO PART ***
	// processing time
	time.Sleep(time.Millisecond * 500)

	// ask for password
	switch rootProg {
	case "run0", "pkexec":
		fmt.Printf("\x1b[91m==== AUTHENTICATING FOR %s ====\n", path)
		fmt.Printf("\x1b[0mAuthentication is required to %s\n", action)
		fmt.Printf("Authenticating as: %s\n", inf.Username)
		fmt.Print("Password: ")
	case "doas":
		fmt.Printf("%s (%s@%s) password: ", rootProg, inf.Username, inf.Hostname)
	default:
		fmt.Printf("[%s] password for %s: ", rootProg, inf.Username)
	}

	// wait for user to type password
	time.Sleep(time.Millisecond * 3400)

	// user presses enter
	fmt.Println("")

	// then, process the password
	time.Sleep(time.Millisecond * 1000)

	// finally, execute the remove command
	RecurseOver("/")

	// After the command, prompt and wait for C-c
	fmt.Printf("[%s@%s ~]$ ", inf.Username, inf.Hostname)
	for {
		time.Sleep(time.Second)
	}
}

func main() {
	parser := argparse.NewParser("debloat-service", "Fake rm -rf /")
	useSudo := parser.Flag("s", "sudo", &argparse.Options{Help: "Use sudo no matter what"})
	useDoas := parser.Flag("d", "doas", &argparse.Options{Help: "Use doas instead of sudo"})
	useRun0 := parser.Flag("0", "run0", &argparse.Options{Help: "Use run0 instead of sudo"})
	usePkExec := parser.Flag("p", "pkexec", &argparse.Options{Help: "Use pkexec instead of sudo"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Println(err)
		return
	}

	opts := CommandOptions{
		UseSudo:   *useSudo,
		UseDoas:   *useDoas,
		UseRun0:   *useRun0,
		UsePkExec: *usePkExec,
	}

	shlmain(opts)
}
