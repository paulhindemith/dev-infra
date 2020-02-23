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
          (mapping)     △    |down
                        |up  ▽
                     __--------__
Step 1              /             \
(Create Kubernetes) \             /
                     ¯¯--------¯¯
          (mapping)     △    |down
                        |up  ▽
                     __--------__
Step 0              /             \
(First Position)    \             /
                     ¯¯--------¯¯
*/

type Element = map[interface{}]interface{}
type mapping = func(Element) (Element, error)

type Mappings interface {
  Name() string
  SimulateUp(Element) (Element, error)
  SimulateDown(Element) (Element, error)
  Up(Element) (Element, error)
  Down(Element) (Element, error)
}

type composite struct {
  simulateUpTo map[int]mapping
  simulateDownTo map[int]mapping
  upTo map[int]mapping
  downTo map[int]mapping
  stepNames map[string]int
  currentStep int
  element Element
}

func Composite(ss ...Mappings) *composite {
  b := &composite{
    stepNames: map[string]int{"0": 0},
    simulateUpTo: map[int]mapping{},
    simulateDownTo: map[int]mapping{},
    upTo: map[int]mapping{},
    downTo: map[int]mapping{},
    currentStep: 0,
    element: Element{},
  }
  for i, s := range ss {
    b.simulateUpTo[i+1] = s.SimulateUp
    b.simulateDownTo[i] = s.SimulateDown
    b.upTo[i+1] = s.Up
    b.downTo[i] = s.Down
    b.stepNames[s.Name()] = i+1
  }
  return b
}
func (c *composite) GetElement() Element {
  return c.element
}
func (c *composite) GetCurrentStep() string {
  for name, i := range c.stepNames {
    if i == c.currentStep {
      return name
    }
  }
  return ""
}

func (c *composite) SimulateAt(name string) error {
  var (
    step int
    ok bool
  )
  if step, ok = c.stepNames[name]; !ok {
    return fmt.Errorf("Got unknown script: %s", name)
  } else if step == c.currentStep {
    return nil
  }

  if step > c.currentStep {
    return c.simulateUp(step)
  }
  return c.simulateDown(step)
}

func (c *composite) ReproduceAt(name string) error {
  var (
    step int
    ok bool
  )
  if step, ok = c.stepNames[name]; !ok {
    return fmt.Errorf("Got unknown script: %s", name)
  } else if step == c.currentStep {
    return nil
  }

  if step > c.currentStep {
    return c.up(step)
  }
  return  c.down(step)

}

func (c *composite) simulateUp(s int) error {
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

func (c *composite) simulateDown(s int) error {
  var (
    e = c.element
    err error
  )
  for i:=c.currentStep-1; i>=s; i-- {
    f := c.simulateDownTo[i]
    if e, err = f(e); err != nil {
      return err
    }
    c.element = e
    c.currentStep = i
  }
  return nil
}

func (c *composite) up(s int) error {
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

func (c *composite) down(s int) error {
  var (
    e = c.element
    err error
  )
  for i:=c.currentStep-1; i>=s; i-- {
    f := c.downTo[i]
    if e, err = f(e); err != nil {
      return err
    }
    c.element = e
    c.currentStep = i
  }
  return nil
}
