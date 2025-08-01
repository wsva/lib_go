//go:build linux
// +build linux

/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kill

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	bscerrors "github.com/wsva/lib_go/kill/errors"

	"k8s.io/klog/v2"
)

// PKill finds process(es) using a regular expression and send a specified
// signal to each process.
func PKill(name string, sig syscall.Signal) error {
	if len(name) == 0 {
		return fmt.Errorf("name should not be empty")
	}
	re, err := regexp.Compile(name)
	if err != nil {
		return err
	}
	pids := GetPids(re)
	if len(pids) == 0 {
		return fmt.Errorf("unable to fetch pids for process name : %q", name)
	}
	errList := []error{}
	for _, pid := range pids {
		if err = syscall.Kill(pid, sig); err != nil {
			errList = append(errList, err)
		}
	}
	return bscerrors.NewAggregate(errList)
}

// PidOf finds process(es) with a specified name (regexp match)
// and return their pid(s).
func PidOf(name string) ([]int, error) {
	if len(name) == 0 {
		return []int{}, fmt.Errorf("name should not be empty")
	}
	re, err := regexp.Compile("(^|/)" + name + "$")
	if err != nil {
		return []int{}, err
	}
	return GetPids(re), nil
}

func GetPids(re *regexp.Regexp) []int {
	pids := []int{}

	dirFD, err := os.Open("/proc")
	if err != nil {
		return nil
	}
	defer dirFD.Close()

	for {
		// Read a small number at a time in case there are many entries, we don't want to
		// allocate a lot here.
		ls, err := dirFD.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}

		for _, entry := range ls {
			if !entry.IsDir() {
				continue
			}

			// If the directory is not a number (i.e. not a PID), skip it
			pid, err := strconv.Atoi(entry.Name())
			if err != nil {
				continue
			}

			cmdline, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "cmdline"))
			if err != nil {
				klog.V(4).Infof("Error reading file %s: %+v", filepath.Join("/proc", entry.Name(), "cmdline"), err)
				continue
			}

			// The bytes we read have '\0' as a separator for the command line
			parts := bytes.SplitN(cmdline, []byte{0}, 2)
			if len(parts) == 0 {
				continue
			}
			// Split the command line itself we are interested in just the first part
			exe := strings.FieldsFunc(string(parts[0]), func(c rune) bool {
				return unicode.IsSpace(c) || c == ':'
			})
			if len(exe) == 0 {
				continue
			}
			// Check if the name of the executable is what we are looking for
			if re.MatchString(exe[0]) {
				// Grab the PID from the directory path
				pids = append(pids, pid)
			}
		}
	}

	return pids
}
