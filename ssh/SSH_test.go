//go:build linux

package ssh_test

import (
	"fmt"
	"testing"

	ssh "github.com/wsva/lib_go/ssh"
)

func TestSSH(T *testing.T) {
	s := ssh.SSH{
		IP:       "127.0.0.1",
		Port:     "22",
		Username: "username",
		Password: "password",
	}

	fmt.Println("exec 01")
	output1, err := s.ExecV01("whoami")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output1))

	fmt.Println("exec 02")
	output1, err = s.ExecV02([]string{"whoami"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output1))

	fmt.Println("exec 03")
	output2, output3, err := s.ExecV03(
		[]ssh.SSHCmd{
			*ssh.NewSSHCmd("whoami"),
		},
		false,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output2)
	fmt.Println(output3)

	fmt.Println("exec 04")
	output2, output3, err = s.ExecV04(
		ssh.NewSSHSuTo("root", "password"),
		[]ssh.SSHCmd{
			*ssh.NewSSHCmd("whoami"),
		},
		false,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output2)
	fmt.Println(output3)
}
