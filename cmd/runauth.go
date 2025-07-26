package cmd

import (
	"fmt"
	"time"
)

func Auth(rootProg string, inf SysInfo) {
	// *** SUDO PART ***
	// processing time
	time.Sleep(time.Millisecond * 500)

	path, action := "", ""

	switch rootProg {
	case "run0":
		path = "org.freedesktop.systemd1.manage-units"
		action = "start transient unit 'run-p989-i990.service'."
	case "pkexec":
		path = "org.freedesktop.policykit.exec"
		action = "run '/bin/rm' as the super user"
	}

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
}
