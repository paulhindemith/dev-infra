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
  "path"
  "strconv"
)

func keys(m map[string]int) []string {
    ks := make([]string, len(m), len(m))
    for k, v := range m {
        ks[v] = k
    }
    return ks
}

func Server(comp Interface) *http.Server {
  mux := http.NewServeMux()
  mux.HandleFunc("/api/scripts", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set( "Access-Control-Allow-Origin", "*")
    bytes, err := json.Marshal(keys(comp.(*composite).stepNames))
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "JSON marshal error: %v", err)
      return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, string(bytes))
  })

  mux.HandleFunc("/api/simulate/", handleFunc(comp))
  mux.HandleFunc("/api/reproduce/", handleFunc(comp))
  s := &http.Server{
    Addr: "127.0.0.1:8003",
    Handler: mux,
  }
  return s
}

func handleFunc(comp Interface) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set( "Access-Control-Allow-Origin", "*")
    var err error
    var stepId int
    if stepId, err = strconv.Atoi(path.Base(r.URL.Path)); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "Could not strconv: %v", err)
    }
    var stepName string
    for k, v := range comp.(*composite).stepNames {
      if v == stepId {
        stepName = k
      }
    }
    if stepName == "" {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "step is not found.")
    }

    switch path.Dir(r.URL.Path) {
    case "/api/simulate":
      err = comp.SimulateAt(stepName)
    case "/api/reproduce":
      err = comp.ReproduceAt(stepName)
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
  }
}
