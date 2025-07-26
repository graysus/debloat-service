package cmd

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

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
