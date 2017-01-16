package skeleton

import (
	"logrus"
	"os"
	"os/exec"
	"syscall"
	"fmt"
)

func Skeleton(args []string) {

	switch args[0] {
		case "run":
			parent(args)
		case "child":
			child(args)
		default:
			logrus.Panicln("what should I do")
	}
}

func parent(args []string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args[1:]...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, }

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%d\n", os.Getpid())

	if err := cmd.Run(); err != nil {
		logrus.Debugln(err)
		os.Exit(1)
	}
}

func child(args []string) {
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%d\n", os.Getpid())

	if err := cmd.Run(); err != nil {
		logrus.Debugln(err)
		os.Exit(1)
	}

}

func must(err error) {
	if err != nil {
		logrus.Panicln(err)
	}
}
