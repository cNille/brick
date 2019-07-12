package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "init":
		initBrickStructure()
	case "run":
		reqInit(parent)
	case "child":
		reqInit(child)
	default:
		panic("what?")

	}
}

type fn func()

func reqInit(callback fn) {
	if checkInit() {
		callback()
		return
	}
	fmt.Println("Brick is not initiated. Pleas run brick init first.")
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// This only works on linux. Is there a hack for mac?
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	// Limit process with controlgroup
	cg()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname([]byte("container"))
	syscall.Chroot(OriginJail)
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	syscall.Mount("dev", "dev", "dev", 0, "")
	defer syscall.Unmount("proc", 0)

	cmd.Run()
}

// Controlgroup, set limits to container
func cg() {

	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")

	// Create a controlgroup-directory called brick
	os.Mkdir(filepath.Join(pids, "brick"), 0755)

	// Set max amounts of processes allowed to 20
	must(ioutil.WriteFile(filepath.Join(pids, "brick/pids.max"), []byte("20"), 0700))

	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "brick/notify_on_release"), []byte("1"), 0700))

	// Set the current process id into the cgroup.procs
	must(ioutil.WriteFile(filepath.Join(pids, "brick/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
