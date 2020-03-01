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
  "reflect"
  "errors"

  "github.com/google/go-cmp/cmp"
  "github.com/google/go-cmp/cmp/cmpopts"
)

var IgnoreUnexportedTypes []interface{}

var alwaysEqual = cmp.Comparer(func(_, _ interface{}) bool { return true })

var opt = cmp.FilterValues(func(x, y interface{}) bool {
    vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
    if !vx.IsValid() || !vy.IsValid() {
      return false
    }
    switch vx.Kind() {
    case reflect.Map:
      return false
    case reflect.Slice:
      return false
    case reflect.Ptr:
      return false
    case reflect.Struct:
      return false
    }
    if vx.Type() != vy.Type() {
      return false
    }
    if vx.Kind() == reflect.Func {
      return x != nil && y != nil
    }
    if x != y {
      zero := reflect.Zero(vx.Type()).Interface()
      return zero != vx.Interface() && zero != vy.Interface()
    }
    return true
}, alwaysEqual)

func Diff(exp, act map[interface{}]interface{}) error {
  if diff := cmp.Diff(exp, act, opt, cmpopts.IgnoreUnexported(IgnoreUnexportedTypes...)); diff != "" {
    return fmt.Errorf("mismatch (-want +got):\n%s", diff)
  }
  return nil
}

func Test(untyped Interface) error {
  var (
    comp *composite
    ok bool
    e_s = map[interface{}]interface{}{}
    e = map[interface{}]interface{}{}
    err error
  )
  if comp, ok = untyped.(*composite); !ok {
    return errors.New("composite type is wrong.")
  }
  last := len(comp.simulateUpTo)
  for i := 1; i < last+1; i++ {
    f_s := comp.simulateUpTo[i]
    if e_s, err = f_s(e_s); err != nil {
      return err
    }
    if i != last {
      if e, err = f_s(e); err != nil {
        return err
      }
      continue
    }
    f := comp.upTo[last]
    if e, err = f(e); err != nil {
      return err
    }
    if err = Diff(e_s, e); err != nil {
      return err
    }
  }

  for i := last-1; i >= 0; i-- {
    f_s := comp.simulateDownTo[i]
    if e_s, err = f_s(e_s); err != nil {
      return err
    }
    if i != last-1 {
      if e, err = f_s(e); err != nil {
        return err
      }
      continue
    }
    f := comp.downTo[last-1]
    if e, err = f(e); err != nil {
      return err
    }
    if err = Diff(e_s, e); err != nil {
      return err
    }
  }
  return nil
}
