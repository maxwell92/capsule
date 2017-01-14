package skeleton

import (
	"logrus"
	"testing"
)

func Test_Skeleton(t *testing.T) {
	/* clear
	args1 := []string{""}
	logrus.Infoln("No args:")
	Skeleton(args1)
	*/

	args2 := []string{"run", "pwd"}
	logrus.Infoln("Args: run")
	Skeleton(args2)
	
	/* clear
	args3 := []string{"child", "pwd"}
	logrus.Infoln("Args: child")
	Skeleton(args3)
	*/
}

