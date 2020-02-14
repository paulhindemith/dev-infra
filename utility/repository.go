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
	"os/exec"
  "strings"
)

var rrd string

func RepoRootDir() (string, error) {
	var (
		out []byte
		err error
	)

	if rrd != "" {
		return rrd, nil
	} else if out, err = exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput(); err != nil {
		return "", err
	}
	rrd = strings.TrimSpace(string(out))
	return rrd, nil
}
