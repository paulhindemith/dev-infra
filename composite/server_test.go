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
  "io/ioutil"
  "net/http"
  "encoding/json"
  "testing"
)

func TestServe(t *testing.T) {
  var testCases = []struct {
    name string
    url string
    expectKey []string
  } {
    {
      name: "simulate A",
      url: "http://localhost:8003?name=A&type=simulate",
      expectKey: []string{"composite.mappingASimulateUp{}"},
    },
    {
      name: "next reproduce B",
      url: "http://localhost:8003?name=B&type=reproduce",
      expectKey: []string{
        "composite.mappingASimulateUp{}",
        "composite.mappingBUp{}",
      },
    },
  }

  r, a, b := &mappingReady{}, &mappingA{}, &mappingB{}
  var (
    comp = Composite(r, a, b)
    resp *http.Response
    body []byte
    element = map[string]string{}
    err error
  )

  go func() {
    if err = Serve(comp); err != nil && err != http.ErrServerClosed {
      t.Fatal(err)
    }
  }()

  for _, tc := range testCases {
    t.Run(tc.name, func (t *testing.T) {
      if resp, err = http.Get(tc.url); err != nil {
        t.Fatal(err)
      }
      if resp.StatusCode != 200 {
        t.Fatal("should be 200")
      }
      defer resp.Body.Close()
      if body, err = ioutil.ReadAll(resp.Body); err != nil {
        t.Fatal(err)
      }
      if err = json.Unmarshal(body,&element); err != nil {
        t.Fatal(err)
      }
      for _, v := range tc.expectKey {
        if element[v] != "true" {
          t.Fatalf("%s must be true. but %s", v, element[v])
        }
      }
    })
  }}
