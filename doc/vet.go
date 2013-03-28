// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package doc

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// This list of deprecated exports is used to find code that has not been
// updated for Go 1.
var deprecatedExports = map[string][]string{
	"bytes":         {"Add"},
	"crypto/aes":    {"Cipher"},
	"crypto/hmac":   {"NewSHA1", "NewSHA256"},
	"crypto/rand":   {"Seed"},
	"encoding/json": {"MarshalForHTML"},
	"encoding/xml":  {"Marshaler", "NewParser", "Parser"},
	"html":          {"NewTokenizer", "Parse"},
	"image":         {"Color", "NRGBAColor", "RGBAColor"},
	"io":            {"Copyn"},
	"log":           {"Exitf"},
	"math":          {"Fabs", "Fmax", "Fmod"},
	"os":            {"Envs", "Error", "Getenverror", "NewError", "Time", "UnixSignal", "Wait"},
	"reflect":       {"MapValue", "Typeof"},
	"runtime":       {"UpdateMemStats"},
	"strconv":       {"Atob", "Atof32", "Atof64", "AtofN", "Atoi64", "Atoui", "Atoui64", "Btoui64", "Ftoa64", "Itoa64", "Uitoa", "Uitoa64"},
	"time":          {"LocalTime", "Nanoseconds", "NanosecondsToLocalTime", "Seconds", "SecondsToLocalTime", "SecondsToUTC"},
	"unicode/utf8":  {"NewString"},
}

type vetVisitor struct {
	importPaths map[string]string
	errors      map[string]token.Pos
}

func (v *vetVisitor) Visit(n ast.Node) ast.Visitor {
	sel, ok := n.(*ast.SelectorExpr)
	if !ok {
		return v
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Obj != nil {
		return v
	}
	importPath := v.importPaths[id.Name]
	for _, name := range deprecatedExports[importPath] {
		if name == sel.Sel.Name {
			v.errors[fmt.Sprintf("%q.%s not found", importPath, sel.Sel.Name)] = n.Pos()
			return v
		}
	}
	return v
}

func (b *builder) vetPackage(pkg *ast.Package) {
	errors := make(map[string]token.Pos)
	for fname, file := range pkg.Files {
		for _, is := range file.Imports {
			importPath, _ := strconv.Unquote(is.Path.Value)
			if !IsValidPath(importPath) &&
				!strings.HasPrefix(importPath, "exp/") &&
				!strings.HasPrefix(importPath, "appengine") {
				errors[fmt.Sprintf("Unrecognized import path %q", importPath)] = is.Pos()
			}
		}
		if importPaths := b.fileImports[fname]; importPaths != nil {
			v := vetVisitor{importPaths: importPaths, errors: errors}
			ast.Walk(&v, file)
		}
	}
	for message, pos := range errors {
		b.pdoc.Errors = append(b.pdoc.Errors,
			fmt.Sprintf("%s (%s)", message, b.fset.Position(pos)))
	}
}
