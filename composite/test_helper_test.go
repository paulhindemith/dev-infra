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
	"time"

	"go.uber.org/zap"
)

type foo struct {
	Name string
}
type fake struct {
	Name string
	Foo string
}
type act struct {
	Name string
	Foo string
}
type ignoreType struct {
	Name string
}

func TestDiff(t *testing.T) {
	logger, _ := zap.NewProduction()

  var testCases = []struct {
    name string
    exp interface{}
    act interface{}
    wantErr bool
  } {
    {
      name: "different type",
      exp: "https://192.168.1.2:44443",
      act: 1,
      wantErr: true,
    },
    {
      name: "different strings",
      exp: "https://192.168.1.2:44443",
      act: "https://localhost",
    },
    {
      name: "empty string",
      exp: "https://192.168.1.2:44443",
      act: "",
      wantErr: true,
    },
    {
      name: "different string in struct",
      exp: struct{Name string}{Name: "localhost"},
      act: struct{Name string}{Name: "google.com"},
    },
		{
      name: "different slice",
      exp: []foo{{Name: "foo"}, {Name: "bar"}},
      act: []foo{{Name: "bar"}},
			wantErr: true,
    },
		{
      name: "different slice with struct will pass because ignoretype works",
      exp: []struct{Name string}{{Name: "foo"}, {Name: "bar"}},
      act: []struct{Name string}{{Name: "bar"}},
    },
    {
      name: "same slice",
      exp: []struct{Name string}{{Name: "foo"}, {Name: "bar"}},
      act: []struct{Name string}{{Name: "bar"}, {Name: "foo"}},
    },
    {
      name: "same function",
      exp: func()error{return nil},
      act: func()error{return nil},
    },
		{
      name: "empty function",
      exp: func()error{return nil},
      act: nil,
      wantErr: true,
    },
		{
      name: "zap logger",
      exp: map[interface{}]interface{}{"string": logger.Sugar()},
      act: map[interface{}]interface{}{"string": logger.Sugar()},
    },
		{
      name: "fake struct1",
      exp: &fake{Name: "foo"},
      act: &act{Foo: "bar"},
    },
		{
      name: "fake struct2",
      exp: &fake{Name: "foo"},
      act: &act{Foo: "bar", Name: ""},
    },
		{
      name: "empty fake struct",
      exp: &fake{Name: "foo"},
      act: nil,
			wantErr: true,
    },
		{
      name: "time",
      exp: time.Now(),
      act: time.Now().Add(-time.Hour),
    },
		{
      name: "time",
      exp: time.Now(),
      act: time.Time{},
			wantErr: true,
    },
		{
      name: "ignore types",
      exp: ignoreType{},
      act: ignoreType{Name: "hoge"},
    },
  }
	RegardAsIdenticalTypesUnlessEitherIsZero(&fake{}, &act{})
	RegardAsIdenticalTypesUnlessEitherIsZero(time.Time{}, time.Time{})
	IgnoreUnexportedTypes(zap.SugaredLogger{})
	IgnoreTypes(ignoreType{})

  for _, tc := range testCases {
    t.Run(tc.name, func (t *testing.T) {
      e_s := map[interface{}]interface{}{}
      e := map[interface{}]interface{}{}
      e_s[struct{}{}] = tc.exp
      e[struct{}{}] = tc.act
			if err := Diff(e_s, e); err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err == nil && tc.wantErr {
				t.Fatal("Coud not catch error")
			}
    })
  }
}

type mapping1 struct {}
func (a *mapping1) Name() string {return "mapping1"}
func (a *mapping1) SimulateUp(e Element) (Element, error) {e[struct{}{}] = true; return e, nil}
func (a *mapping1) SimulateDown(e Element) (Element, error) {e[struct{}{}] = false; return e, nil}
func (a *mapping1) Up(e Element) (Element, error) {e[struct{}{}] = true; return e, nil}
func (a *mapping1) Down(e Element) (Element, error) {e[struct{}{}] = false; return e, nil}

type mapping2 struct {}
func (a *mapping2) Name() string {return "mapping2"}
func (a *mapping2) SimulateUp(e Element) (Element, error) {e[struct{}{}] = "a"; return e, nil}
func (a *mapping2) SimulateDown(e Element) (Element, error) {e[struct{}{}] = "a"; return e, nil}
func (a *mapping2) Up(e Element) (Element, error) {e[struct{}{}] = ""; return e, nil}
func (a *mapping2) Down(e Element) (Element, error) {e[struct{}{}] = "a"; return e, nil}

func TestTest(t *testing.T) {
	var testCases = []struct {
		name string
		mappings []Mappings
		wantErr bool
	} {
		{
			name: "bdd",
			mappings: []Mappings{&mapping1{}},
		},
		{
			name: "error",
			mappings: []Mappings{&mapping1{}, &mapping2{}},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
    t.Run(tc.name, func (t *testing.T) {
			var (
				comp = Composite(tc.mappings...)
			)
			if err := Test(comp); err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err == nil && tc.wantErr {
				t.Fatal("Coud not catch error")
			}
    })
  }
}
