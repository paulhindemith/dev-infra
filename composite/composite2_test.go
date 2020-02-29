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
  "testing"
)

type  mappingA2SimulateUp struct {}
type  mappingA2Up struct {}
type  mappingB2SimulateUp struct {}
type  mappingB2Up struct {}

type mappingA2 struct {}
func (a *mappingA2) Name() string {return "A"}
func (a *mappingA2) SimulateUp(e Element) (Element, error) {e[mappingA2SimulateUp{}] = true; return e, nil}
func (a *mappingA2) Up(e Element) (Element, error) {e[mappingA2Up{}] = true; return e, nil}
type mappingB2 struct {}
func (b *mappingB2) Name() string {return "B"}
func (b *mappingB2) SimulateUp(e Element) (Element, error) {e[mappingB2SimulateUp{}] = true; return e, nil}
func (b *mappingB2) Up(e Element) (Element, error) {e[mappingB2Up{}] = true; return e, nil}
type mappingReady2 struct {}
func (r *mappingReady2) Name() string {return "Ready"}
func (r *mappingReady2) SimulateUp(e Element) (Element, error) {return ready2(e), nil}
func (r *mappingReady2) Up(e Element) (Element, error) {return ready2(e), nil}
func ready2(e Element) Element {
  e[mappingA2SimulateUp{}] = false
  e[mappingA2Up{}] = false
  e[mappingB2SimulateUp{}] = false
  e[mappingB2Up{}] = false
  return e
}
func TestComposite2(t *testing.T) {
  r, a, b := &mappingReady2{}, &mappingA2{}, &mappingB2{}
  var (
    comp = Composite2(r, a, b)
  )

  comp.SimulateAt("A")
  e := comp.GetElement()
  if !e[mappingA2SimulateUp{}].(bool) {
    t.Fatal("mappingA2SimulateUp should be true")
  } else if e[mappingB2SimulateUp{}].(bool) {
    t.Fatal("mappingB2SimulateUp should be false")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }
  comp.SimulateAt("A")
  if !e[mappingA2SimulateUp{}].(bool) {
    t.Fatal("mappingA2SimulateUp should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

  comp.ReproduceAt("B")
  if !e[mappingB2Up{}].(bool) {
    t.Fatal("mappingB2Up should be true")
  }
  if comp.GetCurrentStep() != "B" {
    t.Fatal("CurrentStep should be B")
  }
}
