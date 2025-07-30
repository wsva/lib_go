package cmd

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

func ExecInShell(command string, timeoutSeconds int) (string, error) {
	duration := time.Duration(timeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh")
	input := bytes.NewBuffer(nil)
	output := bytes.NewBuffer(nil)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = output
	input.WriteString(command)
	err := cmd.Run()
	return output.String(), err
}

func ExecInWindowsCmd(command string, timeoutSeconds int) (string, error) {
	duration := time.Duration(timeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	args := []string{"/C"}
	args = append(args, strings.Split(command, " ")...)
	output, err := exec.CommandContext(ctx, "cmd.exe", args...).Output()
	return string(output), err
}

func ExecInWindowsCmd1(command string, timeoutSeconds int) (string, error) {
	duration := time.Duration(timeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cmd.exe")
	input := bytes.NewBuffer(nil)
	output := bytes.NewBuffer(nil)
	cmd.Stdin = input
	cmd.Stdout = output
	if command[len(command)-1] != '\n' {
		command += "\n"
	}
	input.WriteString(command)
	err := cmd.Run()
	return mahonia.NewDecoder("gbk").ConvertString(output.String()), err
}
