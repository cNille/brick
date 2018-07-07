package main

import (
	"fmt"
	//"io/ioutil"
	"os"
	"os/exec"
	//"path/filepath"
	//"strconv"
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

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot(OriginJail))
	must(syscall.Chdir("/"))
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))
	defer must(syscall.Unmount("proc", 0))

	must(cmd.Run())

}

// func cg() {
// 	cgroups := "/sys/fs/cgroup/"
// 	pids := filepath.Join(cgroups, "pids")
// 	os.Mkdir(filepath.Join(pids, "nille"), 0755)
// 	must(ioutil.WriteFile(filepath.Join(pids, "nille/pids.max"), []byte("20"), 0700))
// 	// Removes the new cgroup in place after the container exits
// 	must(ioutil.WriteFile(filepath.Join(pids, "nille/notify_on_release"), []byte("1"), 0700))
// 	must(ioutil.WriteFile(filepath.Join(pids, "nille/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
// }

func must(err error) {
	if err != nil {
		panic(err)
	}
}
