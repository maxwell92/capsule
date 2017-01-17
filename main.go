package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("what should I do")
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	/*
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	*/

	cmd.SysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET,
}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("pid: %d\n", os.Getpid())
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Printf("ERROR", err)
		os.Exit(1)
	}
	fmt.Printf("hostname: %s\n", hostName)

	if err := cmd.Run(); err != nil {
		fmt.Printf("ERROR", err)
		os.Exit(1)
	}
}

func child() {
	/*
	err := syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, "")
	if err != nil {
		fmt.Printf("Mount rootfs\n")
		panic(err)
	}

	err = os.MkdirAll("rootfs/oldrootfs", 0700)
	if err != nil {
		fmt.Printf("MkdirAll rootfs\n")
		panic(err)
	}

	err = syscall.PivotRoot("rootfs", "rootfs/oldrootfs")
	if err != nil {
		fmt.Printf("PivotRoot rootfs\n")
		panic(err)
	}

	err = os.Chdir("/")
	if err != nil {
		fmt.Printf("Chdir rootfs\n")
		panic(err)
	}
	*/

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("pid: %d\n", os.Getpid())

	syscall.Sethostname([]byte("NewNamespace"))

	hostName, err := os.Hostname()
	if err != nil {
		fmt.Printf("ERROR", err)
		os.Exit(1)
	}
	fmt.Printf("hostname: %s\n", hostName)

	if err := cmd.Run(); err != nil {
		fmt.Printf("ERROR", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
