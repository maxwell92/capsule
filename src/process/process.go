package process

import (
	"fmt"
	"logrus"
	"os"
	"syscall"
)

func Create() error {

	spAttr := &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}

	pAttr := &os.ProcAttr{
		Sys: spAttr,
	}

	// it's says no such file or directory
	np, err := os.StartProcess("new process", nil, pAttr)
	if err != nil {
		logrus.Println(err)
		return err
	}

	fmt.Printf("New Process Pid: %d\n", np.Pid)
	fmt.Printf("Pid: %d\n", os.Getpid())
	return nil
}
