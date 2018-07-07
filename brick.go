package main

import (
	"fmt"
	"github.com/cNille/brick/util"
	"os"
	"os/exec"
)

var (
	workDirectory        = "/vagrant"
	originRootFileSystem = "/vagrant/data/rootfs.tar.gz"
	nodetar              = "/vagrant/data/node-v8.11.3-linux-x64.tar.xz"
	brickHome            = "/srv/brick"
	nodejs               = "/srv/brick/rootfs/usr/local/lib"
	OriginJail           = "/srv/brick/rootfs"
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
	util.EnsureDir(brickHome)

	// Create a dir for a jail with rootfilesystem.
	// A process can be isolated in the folder with Chroot.
	if _, err := os.Stat(OriginJail); os.IsNotExist(err) {
		fmt.Println("Unzip ubuntu")
		cmd := exec.Command("tar", "-zxf", originRootFileSystem, "--directory", brickHome)
		cmd.Run()
	}

	err := util.Copy(workDirectory+"/data/install_node.sh", OriginJail+"/install_node.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	// install node
	fmt.Println("Unzip nodejs")
	cmd := exec.Command("tar", "-xf", nodetar, "--directory", nodejs)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Unable to unzip nodejs...")
		fmt.Println(err)
		return
	}

	// mv node folder
	cmd = exec.Command("mv", nodejs+"/node-v8.11.3-linux-x64", nodejs+"/node-v8.11.3")
	cmd.Run()

	// mv profile with node export
	cmd = exec.Command("cp", workDirectory+"/data/.profile", OriginJail)
	cmd.Run()

	fmt.Printf("Done!\n")
}
