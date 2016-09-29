package process

import (
	"fmt"
	"testing"
)

func Test_Create(t *testing.T) {
	if err := Create(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
}
