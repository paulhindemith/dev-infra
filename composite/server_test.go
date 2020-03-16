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
  "context"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "testing"
)

func TestGetScript(t *testing.T) {
  url := "http://localhost:8003/api/scripts"

  r, a, b := &mappingReady{}, &mappingA{}, &mappingB{}
  var (
    comp = Composite(r, a, b)
    resp *http.Response
    body []byte
    scripts = []string{}
    s = Server(comp)
    err error
  )

  go func() {
    if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      t.Fatal(err)
    }
  }()
  defer s.Shutdown(context.Background())

  if resp, err = http.Get(url); err != nil {
    t.Fatal(err)
  }
  if resp.StatusCode != 200 {
    t.Fatalf("should be 200. but %d", resp.StatusCode)
  }
  defer resp.Body.Close()
  if body, err = ioutil.ReadAll(resp.Body); err != nil {
    t.Fatal(err)
  }
  if err = json.Unmarshal(body, &scripts); err != nil {
    t.Fatal(err)
  }
  if scripts[0] != "0" {
    t.Fatalf(`scripts[0] must be "0", but %s`, scripts[0])
  }
  if scripts[1] != "Ready" {
    t.Fatalf(`scripts[1] must be "Ready", but %s`, scripts[0])
  }
  if scripts[2] != "A" {
    t.Fatalf(`scripts[2] must be "A", but %s`, scripts[0])
  }
  if scripts[3] != "B" {
    t.Fatalf(`scripts[3] must be "B", but %s`, scripts[0])
  }
  if len(scripts) != 4 {
    t.Fatalf("scripts length must be 3, but %d", len(scripts))
  }
}

func TestRunScript(t *testing.T) {
  var testCases = []struct {
    name string
    url string
    expectKey []string
  } {
    {
      name: "simulate 1",
      url: "http://localhost:8003/api/simulate/1",
      expectKey: []string{},
    },
    {
      name: "next reproduce 2",
      url: "http://localhost:8003/api/simulate/2",
      expectKey: []string{
        "composite.mappingASimulateUp{}",
      },
    },
    {
      name: "next reproduce 3",
      url: "http://localhost:8003/api/reproduce/3",
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
    s = Server(comp);
    err error
  )

  go func() {
    if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      t.Fatal(err)
    }
  }()
  defer s.Shutdown(context.Background())


  for _, tc := range testCases {
    t.Run(tc.name, func (t *testing.T) {
      if resp, err = http.Get(tc.url); err != nil {
        t.Fatal(err)
      }
      if resp.StatusCode != 200 {
        t.Fatalf("should be 200. but %d", resp.StatusCode)
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
