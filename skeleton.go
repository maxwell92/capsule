package skeleton

import (
	"logrus"
	"os"
	"os/exec"
	"syscall"
)

func Skeleton() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		logrus.Panicln("what should I do")
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Debugln(err)
		os.Exit(1)
	}
}

func child() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
