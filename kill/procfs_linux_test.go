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

package kill_test

import (
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wsva/lib_go/kill"
)

func TestPidOf(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skipf("not supported on GOOS=%s", runtime.GOOS)
	}
	pids, err := kill.PidOf(filepath.Base(os.Args[0]))
	assert.Empty(t, err)
	assert.NotZero(t, pids)
	assert.Contains(t, pids, os.Getpid())
}

func TestPKill(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skipf("not supported on GOOS=%s", runtime.GOOS)
	}
	sig := syscall.SIGCONT
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig)
	defer signal.Stop(c)
	kill.PKill(os.Args[0], sig)
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("signal was %v, want %v", s, sig)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for %v", sig)
	}
}

func BenchmarkGetPids(b *testing.B) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		b.Skipf("not supported on GOOS=%s", runtime.GOOS)
	}

	re, err := regexp.Compile("(^|/)" + filepath.Base(os.Args[0]) + "$")
	assert.Empty(b, err)

	for i := 0; i < b.N; i++ {
		pids := kill.GetPids(re)

		b.StopTimer()
		assert.NotZero(b, pids)
		assert.Contains(b, pids, os.Getpid())
		b.StartTimer()
	}
}
