/*
Copyright 2024 The GoPlus Authors (goplus.org)
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

package hdqtest

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/hdq"
)

// func(doc hdq.NodeSet) any
type Converter = any

// FromDir tests all html files in a directory.
func FromDir(t *testing.T, sel, relDir string, conv Converter) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	vConv := reflect.ValueOf(conv)
	for _, fi := range fis {
		name := fi.Name()
		if !fi.IsDir() || strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, sel, vConv)
		})
	}
}

func testFrom(t *testing.T, pkgDir, sel string, conv reflect.Value) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	log.Println("Parsing", pkgDir)
	in := pkgDir + "/in.html"
	out := pkgDir + "/out.json"
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal("ReadFile failed:", err)
	}
	expected := string(b)
	ret := ConvFile(in, conv)
	retb, _ := json.MarshalIndent(ret, "", "\t")
	if v := string(retb); v != expected {
		t.Fatalf("\n==> got:\n%s\n==> expected:\n%s\n", v, expected)
	}
}

// ConvFile converts a html source to an object.
func ConvFile(in any, conv reflect.Value) any {
	doc := reflect.ValueOf(hdq.Source(in))
	out := conv.Call([]reflect.Value{doc})
	return out[0].Interface()
}
