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
)
/*
                     __--------__
Step 2              /             \
(Apply Knative)     \             /
                     ¯¯--------¯¯
          (mapping)     △
                        |up
                     __--------__
Step 1              /             \
(Create Kubernetes) \             /
                     ¯¯--------¯¯
          (mapping)     △
                        |up
                     __--------__
Step 0              /             \
(First Position)    \             /
                     ¯¯--------¯¯
*/

type Mappings2 interface {
  Name() string
  SimulateUp(Element) (Element, error)
  Up(Element) (Element, error)
}

type composite2 struct {
  simulateUpTo map[int]mapping
  upTo map[int]mapping
  stepNames map[string]int
  currentStep int
  element Element
}

func Composite2(ss ...Mappings2) Interface {
  b := &composite{
    stepNames: map[string]int{"0": 0},
    simulateUpTo: map[int]mapping{},
    upTo: map[int]mapping{},
    currentStep: 0,
    element: Element{},
  }
  for i, s := range ss {
    b.simulateUpTo[i+1] = s.SimulateUp
    b.upTo[i+1] = s.Up
    b.stepNames[s.Name()] = i+1
  }
  return b
}
func (c *composite2) GetElement() Element {
  return c.element
}
func (c *composite2) GetCurrentStep() string {
  for name, i := range c.stepNames {
    if i == c.currentStep {
      return name
    }
  }
  return ""
}

func (c *composite2) SimulateAt(name string) error {
  var (
    step int
    ok bool
  )
  if step, ok = c.stepNames[name]; !ok {
    return fmt.Errorf("Got unknown script: %s", name)
  } else if step == c.currentStep {
    return nil
  }

  if step < c.currentStep {
    return fmt.Errorf("Can not simulateDown to %s", name)
  }
  return c.simulateUp(step)
}

func (c *composite2) ReproduceAt(name string) error {
  var (
    step int
    ok bool
  )
  if step, ok = c.stepNames[name]; !ok {
    return fmt.Errorf("Got unknown script: %s", name)
  } else if step == c.currentStep {
    return nil
  }

  if step < c.currentStep {
    return fmt.Errorf("Can not Down to %s", name)
  }
  return  c.up(step)

}

func (c *composite2) simulateUp(s int) error {
  var (
    e = c.element
    err error
  )
  for i:=c.currentStep+1; i<=s; i++ {
    f := c.simulateUpTo[i]
    if e, err = f(e); err != nil {
      return err
    }
    c.element = e
    c.currentStep = i
  }
  return nil
}

func (c *composite2) up(s int) error {
  var (
    e = c.element
    err error
  )
  for i:=c.currentStep+1; i<=s; i++ {
    f := c.upTo[i]
    if e, err = f(e); err != nil {
      return err
    }
    c.element = e
    c.currentStep = i
  }
  return nil
}
