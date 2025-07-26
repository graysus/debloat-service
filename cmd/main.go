package cmd

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type CommandOptions struct {
	UseSudo, UseDoas, UseRun0, UsePkExec bool
}

func Main() {
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
