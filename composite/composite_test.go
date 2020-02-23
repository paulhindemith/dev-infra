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

type  mappingASimulateUp struct {}
type  mappingASimulateDown struct {}
type  mappingAUp struct {}
type  mappingADown struct {}
type  mappingBSimulateUp struct {}
type  mappingBSimulateDown struct {}
type  mappingBUp struct {}
type  mappingBDown struct {}

type mappingA struct {}
func (a *mappingA) Name() string {return "A"}
func (a *mappingA) SimulateUp(e Element) (Element, error) {e[mappingASimulateUp{}] = true; return e, nil}
func (a *mappingA) SimulateDown(e Element) (Element, error) {e[mappingASimulateDown{}] = true; return e, nil}
func (a *mappingA) Up(e Element) (Element, error) {e[mappingAUp{}] = true; return e, nil}
func (a *mappingA) Down(e Element) (Element, error) {e[mappingADown{}] = true; return e, nil}
type mappingB struct {}
func (b *mappingB) Name() string {return "B"}
func (b *mappingB) SimulateUp(e Element) (Element, error) {e[mappingBSimulateUp{}] = true; return e, nil}
func (b *mappingB) SimulateDown(e Element) (Element, error) {e[mappingBSimulateDown{}] = true; return e, nil}
func (b *mappingB) Up(e Element) (Element, error) {e[mappingBUp{}] = true; return e, nil}
func (b *mappingB) Down(e Element) (Element, error) {e[mappingBDown{}] = true; return e, nil}
type mappingReady struct {}
func (r *mappingReady) Name() string {return "Ready"}
func (r *mappingReady) SimulateUp(e Element) (Element, error) {return ready(e), nil}
func (r *mappingReady) SimulateDown(e Element) (Element, error) {return ready(e), nil}
func (r *mappingReady) Up(e Element) (Element, error) {return ready(e), nil}
func (r *mappingReady) Down(e Element) (Element, error) {return ready(e), nil}
func ready(e Element) Element {
  e[mappingASimulateUp{}] = false
  e[mappingASimulateDown{}] = false
  e[mappingAUp{}] = false
  e[mappingADown{}] = false
  e[mappingBSimulateUp{}] = false
  e[mappingBSimulateDown{}] = false
  e[mappingBUp{}] = false
  e[mappingBDown{}] = false
  return e
}
func TestComposite(t *testing.T) {
  r, a, b := &mappingReady{}, &mappingA{}, &mappingB{}
  var (
    comp = Composite(r, a, b)
  )

  comp.SimulateAt("A")
  e := comp.GetElement()
  if !e[mappingASimulateUp{}].(bool) {
    t.Fatal("mappingASimulateUp should be true")
  } else if e[mappingBSimulateUp{}].(bool) {
    t.Fatal("mappingBSimulateUp should be false")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }
  comp.SimulateAt("A")
  if !e[mappingASimulateUp{}].(bool) {
    t.Fatal("mappingASimulateUp should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

  comp.SimulateAt("B")
  if !e[mappingBSimulateUp{}].(bool) {
    t.Fatal("mappingBSimulateUp should be true")
  }
  if comp.GetCurrentStep() != "B" {
    t.Fatal("CurrentStep should be B")
  }

  comp.SimulateAt("A")
  if !e[mappingBSimulateDown{}].(bool) {
    t.Fatal("mappingBSimulateDown should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }


  comp.ReproduceAt("A")
  if e[mappingBUp{}].(bool) {
    t.Fatal("mappingBUp should be false")
  }
  if e[mappingBDown{}].(bool) {
    t.Fatal("mappingBDown should be false")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

  comp.ReproduceAt("B")
  if !e[mappingBUp{}].(bool) {
    t.Fatal("mappingBUp should be true")
  }
  if comp.GetCurrentStep() != "B" {
    t.Fatal("CurrentStep should be B")
  }

  comp.ReproduceAt("A")
  if !e[mappingBDown{}].(bool) {
    t.Fatal("mappingBDown should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

}
