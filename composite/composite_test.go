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
  "testing"
)

type mappingA struct {}
type mappingASimulateUpKey struct{}
type mappingASimulateDownKey struct{}
type mappingAUpKey struct{}
type mappingADownKey struct{}
func (a *mappingA) Name() string {return "A"}
func (a *mappingA) SimulateUp(c context.Context) (context.Context, error) {return context.WithValue(c, mappingASimulateUpKey{}, true), nil}
func (a *mappingA) SimulateDown(c context.Context) (context.Context, error) {return context.WithValue(c, mappingASimulateDownKey{}, true), nil}
func (a *mappingA) Up(c context.Context) (context.Context, error) {return context.WithValue(c, mappingAUpKey{}, true), nil}
func (a *mappingA) Down(c context.Context) (context.Context, error) {return context.WithValue(c, mappingADownKey{}, true), nil}
type mappingB struct {}
type mappingBSimulateUpKey struct{}
type mappingBSimulateDownKey struct{}
type mappingBUpKey struct{}
type mappingBDownKey struct{}
func (b *mappingB) Name() string {return "B"}
func (b *mappingB) SimulateUp(c context.Context) (context.Context, error) {return context.WithValue(c, mappingBSimulateUpKey{}, true), nil}
func (b *mappingB) SimulateDown(c context.Context) (context.Context, error) {return context.WithValue(c, mappingBSimulateDownKey{}, true), nil}
func (b *mappingB) Up(c context.Context) (context.Context, error) {return context.WithValue(c, mappingBUpKey{}, true), nil}
func (b *mappingB) Down(c context.Context) (context.Context, error) {
  return context.WithValue(c, mappingBDownKey{}, true), nil}
type mappingReady struct {}
func (r *mappingReady) Name() string {return "Ready"}
func (r *mappingReady) SimulateUp(c context.Context) (context.Context, error) {return ready(c), nil}
func (r *mappingReady) SimulateDown(c context.Context) (context.Context, error) {return ready(c), nil}
func (r *mappingReady) Up(c context.Context) (context.Context, error) {return ready(c), nil}
func (r *mappingReady) Down(c context.Context) (context.Context, error) {return ready(c), nil}
func ready(c context.Context) context.Context {
  c = context.WithValue(c, mappingASimulateUpKey{}, false)
  c = context.WithValue(c, mappingASimulateDownKey{}, false)
  c = context.WithValue(c, mappingAUpKey{}, false)
  c = context.WithValue(c, mappingADownKey{}, false)
  c = context.WithValue(c, mappingBSimulateUpKey{}, false)
  c = context.WithValue(c, mappingBSimulateDownKey{}, false)
  c = context.WithValue(c, mappingBUpKey{}, false)
  c = context.WithValue(c, mappingBDownKey{}, false)
  return c
}
func TestComposite(t *testing.T) {
  r, a, b := &mappingReady{}, &mappingA{}, &mappingB{}
  var (
    ctx context.Context
    comp = Composite(r, a, b)
  )

  comp.SimulateAt("A")
  ctx = comp.GetContext()
  if !ctx.Value(mappingASimulateUpKey{}).(bool) {
    t.Fatal("mappingASimulateUpKey should be true")
  } else if ctx.Value(mappingBSimulateUpKey{}).(bool) {
    t.Fatal("mappingBSimulateUpKey should be false")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }
  comp.SimulateAt("A")
  ctx = comp.GetContext()
  if !ctx.Value(mappingASimulateUpKey{}).(bool) {
    t.Fatal("mappingASimulateUpKey should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

  comp.SimulateAt("B")
  ctx = comp.GetContext()
  if !ctx.Value(mappingBSimulateUpKey{}).(bool) {
    t.Fatal("mappingBSimulateUpKey should be true")
  }
  if comp.GetCurrentStep() != "B" {
    t.Fatal("CurrentStep should be B")
  }

  comp.SimulateAt("A")
  ctx = comp.GetContext()
  if !ctx.Value(mappingBSimulateDownKey{}).(bool) {
    t.Fatal("mappingBSimulateDownKey should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }


  comp.ReproduceAt("A")
  ctx = comp.GetContext()
  if ctx.Value(mappingBUpKey{}).(bool) {
    t.Fatal("mappingBUpKey should be false")
  }
  if ctx.Value(mappingBDownKey{}).(bool) {
    t.Fatal("mappingBDownKey should be false")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

  comp.ReproduceAt("B")
  ctx = comp.GetContext()
  if !ctx.Value(mappingBUpKey{}).(bool) {
    t.Fatal("mappingBUpKey should be true")
  }
  if comp.GetCurrentStep() != "B" {
    t.Fatal("CurrentStep should be B")
  }

  comp.ReproduceAt("A")
  ctx = comp.GetContext()
  if !ctx.Value(mappingBDownKey{}).(bool) {
    t.Fatal("mappingBDownKey should be true")
  }
  if comp.GetCurrentStep() != "A" {
    t.Fatal("CurrentStep should be A")
  }

}
