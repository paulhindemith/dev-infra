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

package composite

import (
  "fmt"
  "net/http"
  "encoding/json"
)

func Serve(comp Interface) error {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    var (
      script string
      scriptType string
      err error
    )
    for k, v := range q {
      if k == "name" {
        script = v[0]
      }
      if k == "type" {
        scriptType = v[0]
      }
    }
    switch scriptType {
    case "simulate":
      err = comp.SimulateAt(script)
    case "reproduce":
      err = comp.ReproduceAt(script)
    default:
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "script is not found.")
    }
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "Could not successfully finish script: %v", err)
      return
    }
    e := comp.GetElement()
    ret := map[string]string{}
    for k, v := range e {
      ret[fmt.Sprintf("%#v", k)] = fmt.Sprintf("%#v", v)
    }
    bytes, err := json.Marshal(ret)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "JSON marshal error: %v", err)
      return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, string(bytes))
  })
  return http.ListenAndServe("127.0.0.1:8003", mux)
}
