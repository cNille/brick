package main

import (
	"fmt"
	"github.com/cNille/brick/util"
	"os"
)

var (
	workDirectory = "/srv/brick/"
	OriginJail    = workDirectory + "/jail"
	bin           = "/bin"
	lib           = "/lib"
	lib64         = "/lib64"
	proc          = "/proc"
)

func checkInit() bool {
	if _, err := os.Stat(workDirectory); os.IsNotExist(err) {
		return false
	}
	return true
}

func initBrickStructure() {
	fmt.Printf("Start init...\n")

	// Create a work directory for brick
	util.EnsureDir(workDirectory)

	// Create a dir for the original jail
	// A jail is a folder with own binaries and lib.
	// A process can be isolated in the folder with Chroot.
	util.EnsureDir(OriginJail)
	util.EnsureDir(OriginJail + proc)

	// Copy necessary binaries for the jail-template.
	err := util.Copy(bin, OriginJail+bin)
	if err != nil {
		fmt.Printf("ERROR: unable to create bin. \n")
		fmt.Println(err)
		return
	}

	err = util.Copy(lib, OriginJail+lib)
	if err != nil {
		fmt.Printf("ERROR: unable to create lib \n")
		fmt.Println(err)
		return
	}

	err = util.Copy(lib64, OriginJail+lib64)
	if err != nil {
		fmt.Printf("ERROR: unable to create lib64 \n")
		fmt.Println(err)
		return
	}

	fmt.Printf("Done!")
}

func createJail() {

}
