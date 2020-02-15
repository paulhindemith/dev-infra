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
  "path"
  "testing"
  "os/exec"
)

func TestIsTesting(t *testing.T) {
  if !IsTesting() {
    t.Fatal("must be true")
  }
  if rrd, err := RepoRootDir(); err != nil {
    t.Fatal(err)
  } else if out, err := exec.Command("go", "run", path.Join(rrd, "utility_main/is_testing_main.go")).CombinedOutput(); err != nil {
    t.Fatal(err)
  } else if string(out) == "true" {
    t.Fatal("must false")
  }
}
