/*
Copyright 2020 Paulhindemith

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

package utility

import (
  "testing"
  "os/exec"
  "strings"
)

func getOS() string {
  out, _ := exec.Command("uname").CombinedOutput()
  return strings.TrimSpace(string(out))
}

func TestLinux(t *testing.T) {
  switch getOS() {
  case "Linux":
    if !IsLinux() {
      t.Fatal("must true")
    }
  case "Darwin":
    if IsLinux() {
      t.Fatal("must false")
    }
  default:
    t.Fatal("unknown os")
  }
}

func TestNotLinuxSkipWhenShort(t *testing.T) {
  if testing.Short() {
    t.Fatal("You can not run short test.")
  }
  switch getOS() {
  case "Linux":
    NotLinuxSkipWhenShort(t)
  }
}

func TestDarwin(t *testing.T) {
  switch getOS() {
  case "Linux":
    if IsDarwin() {
      t.Fatal("must false")
    }
  case "Darwin":
    if !IsDarwin() {
      t.Fatal("must true")
    }
  default:
    t.Fatal("unknown os")
  }
}

func TestNotDarwinSkipWhenShort(t *testing.T) {
  if testing.Short() {
    t.Fatal("You can not run short test.")
  }
  switch getOS() {
  case "Darwin":
    NotDarwinSkipWhenShort(t)
  }
}
