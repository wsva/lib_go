//go:build linux

package ssh

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
)

const (
	SSHTimeout = 60 * time.Second
)

type SSHConfig struct {
	TimeoutSSH   time.Duration `json:"TimeoutSSH"`
	TimeoutSpawn time.Duration `json:"TimeoutSpawn"`

	RegPrompt   string `json:"RegPrompt"`
	RegPassword string `json:"RegPassword"`
}

func DefaultSSHConfig() *SSHConfig {
	return &SSHConfig{
		RegPrompt:    `(\$|#)\s*$`,
		RegPassword:  `Password`,
		TimeoutSSH:   SSHTimeout,
		TimeoutSpawn: SSHTimeout,
	}
}

type SSHCmd struct {
	Cmd     string        `json:"Cmd"`
	Timeout time.Duration `json:"Timeout"`
	StdOut  bool          `json:"StdOut"`
	StdErr  bool          `json:"StdErr"`
	Combine bool          `json:"Combine"`
}

func NewSSHCmd(cmd string) *SSHCmd {
	return &SSHCmd{
		Cmd:     cmd,
		Timeout: SSHTimeout,
		StdOut:  true,
		StdErr:  true,
		Combine: true,
	}
}

type SSHSuTo struct {
	Username       string        `json:"Username"`
	Password       string        `json:"Password"`
	ExpectPassword bool          `json:"ExpectPassword"`
	Timeout        time.Duration `json:"Timeout"`
}

func NewSSHSuTo(username, password string) *SSHSuTo {
	return &SSHSuTo{
		Username:       username,
		Password:       password,
		ExpectPassword: true,
		Timeout:        SSHTimeout,
	}
}

type SSH struct {
	IP        string      `json:"IP"`
	Port      string      `json:"Port"`
	Username  string      `json:"Username"`
	Password  string      `json:"Password"`
	SSHConfig *SSHConfig  `json:"-"`
	Client    *ssh.Client `json:"-"`
}

func (s *SSH) Dial() error {
	if s.SSHConfig == nil {
		s.SSHConfig = DefaultSSHConfig()
	}
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         s.SSHConfig.TimeoutSSH,
	}
	client, err := ssh.Dial("tcp", s.IP+":"+s.Port, config)
	if err != nil {
		return err
	}
	s.Client = client
	return nil
}

func (s *SSH) Close() error {
	return s.Client.Close()
}

/*
最简单的，执行一个命令
1，没有Timeout
*/
func (s *SSH) ExecV01(cmd string) (string, error) {
	if s.Client == nil {
		err := s.Dial()
		if err != nil {
			return "", err
		}
	}
	session, err := s.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create session error: %w", err)
	}
	defer session.Close()
	buf, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("run command error: %w", err)
	}
	return string(buf), nil
}

/*
执行多个命令
实现方式：输入输出重定向
1，没有Timeout
*/
func (s *SSH) ExecV02(cmdList []string) (string, error) {
	if s.Client == nil {
		err := s.Dial()
		if err != nil {
			return "", err
		}
	}
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	input := bytes.NewBuffer(nil)
	output := bytes.NewBuffer(nil)

	session.Stdin = input
	session.Stdout = output
	session.Stderr = output

	if err := session.Shell(); err != nil {
		return "", err
	}

	for _, v := range cmdList {
		input.WriteString(s.appendCharReturn(v))
	}
	input.WriteString("exit\n")

	err = session.Wait()
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

/*
执行多个命令
实现方式：expect
1，支持给每个命令指定超时时间
2，仅命令执行的输出，包括其他输出的全部输出
3，支持设置中途报错是否终止执行
*/
func (s *SSH) ExecV03(cmdList []SSHCmd, stop bool) (string, string, error) {
	if s.Client == nil {
		err := s.Dial()
		if err != nil {
			return "", "", err
		}
	}

	e, _, err := expect.SpawnSSH(s.Client, SSHTimeout)
	if err != nil {
		return "", "", err
	}
	defer e.Close()

	var result, resultAll strings.Builder
	regPrompt := regexp.MustCompile(s.SSHConfig.RegPrompt)

	output, _, err := e.Expect(regPrompt, s.SSHConfig.TimeoutSpawn)
	resultAll.WriteString(s.processOutput(output))
	if err != nil {
		return result.String(), resultAll.String(), err
	}

	for _, v := range cmdList {
		if v.Timeout == 0 {
			v.Timeout = SSHTimeout
		}
		err = e.Send(s.appendCharReturn(v.Cmd))
		if err != nil {
			if stop {
				return result.String(), resultAll.String(), err
			}
			fmt.Println(err)
		}
		output, _, err = e.Expect(regPrompt, v.Timeout)
		resultAll.WriteString(s.processOutput(output))
		if v.StdOut || v.StdErr {
			result.WriteString(s.processOutput(s.deleteLastLine(output)))
		}
		if err != nil {
			if stop {
				return result.String(), resultAll.String(), err
			}
			fmt.Println(err)
		}
	}

	return result.String(), resultAll.String(), nil
}

/*
su之后再执行多个命令
实现方式：expect
1，支持给每个命令指定超时时间
2，仅命令执行的输出，包括其他输出的全部输出
3，支持设置中途报错是否终止执行
*/
func (s *SSH) ExecV04(su *SSHSuTo, cmdList []SSHCmd, stop bool) (string, string, error) {
	if s.Client == nil {
		err := s.Dial()
		if err != nil {
			return "", "", err
		}
	}

	e, _, err := expect.SpawnSSH(s.Client, SSHTimeout)
	if err != nil {
		return "", "", err
	}
	defer e.Close()

	var result, resultAll strings.Builder
	regPrompt := regexp.MustCompile(s.SSHConfig.RegPrompt)
	regPassword := regexp.MustCompile(s.SSHConfig.RegPassword)

	output, _, err := e.Expect(regPrompt, s.SSHConfig.TimeoutSpawn)
	resultAll.WriteString(s.processOutput(output))
	if err != nil {
		return result.String(), resultAll.String(), err
	}

	//LANG=en 防止出现非Password的提示
	err = e.Send(fmt.Sprintf("LANG=en su - %s\n", su.Username))
	if err != nil {
		return result.String(), resultAll.String(), err
	}
	if su.ExpectPassword {
		output, _, err = e.Expect(regPassword, su.Timeout)
		resultAll.WriteString(s.processOutput(output))
		if err != nil {
			return result.String(), resultAll.String(), err
		}
		err = e.Send(su.Password + "\n")
		if err != nil {
			return result.String(), resultAll.String(), err
		}
	}
	output, _, err = e.Expect(regPrompt, su.Timeout)
	resultAll.WriteString(s.processOutput(output))
	if err != nil {
		return result.String(), resultAll.String(), err
	}

	for _, v := range cmdList {
		if v.Timeout == 0 {
			v.Timeout = SSHTimeout
		}
		err = e.Send(s.appendCharReturn(v.Cmd))
		if err != nil {
			if stop {
				return result.String(), resultAll.String(), err
			}
			fmt.Println(err)
		}
		output, _, err = e.Expect(regPrompt, v.Timeout)
		resultAll.WriteString(s.processOutput(output))
		if v.StdOut || v.StdErr {
			result.WriteString(s.processOutput(s.deleteLastLine(output)))
		}
		if err != nil {
			if stop {
				return result.String(), resultAll.String(), err
			}
			fmt.Println(err)
		}
	}

	return result.String(), resultAll.String(), nil
}

/*
problem: there is [27 91 63 50 48 48 52 108 13] in every output of goexpect
*/
func (s *SSH) processOutput(output string) string {
	pattern := []byte{27, 91, 63, 50, 48, 48, 52, 108, 13}
	return strings.ReplaceAll(output, string(pattern), "")
}

func (s *SSH) appendCharReturn(content string) string {
	if len(content) > 0 && content[len(content)-1] == '\n' {
		return content
	}
	return content + "\n"
}

func (s *SSH) deleteLastLine(content string) string {
	result := ""
	lines := strings.Split(content, "\n")
	count := len(lines)
	for i := 0; i < count-1; i++ {
		if i < count-2 {
			result += lines[i] + "\n"
		} else {
			result += lines[i]
		}
	}
	return result
}
