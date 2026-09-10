//go:build ignore

// Extracts the Go SDK's public service surface into the shared api-surface JSON.
//
// Uses go/ast from the standard library rather than go/packages: parsing the service
// package's files needs no module loading, so this runs on any toolchain.
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type method struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
	Doc       string `json:"doc"`
}

type service struct {
	Name    string   `json:"name"`
	Methods []method `json:"methods"`
}

type surface struct {
	SDK      string    `json:"sdk"`
	Source   string    `json:"source"`
	Services []service `json:"services"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: extract-go.go <service-dir> <out.json>")
		os.Exit(2)
	}
	dir, out := os.Args[1], os.Args[2]

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", dir, err)
		os.Exit(1)
	}

	byService := map[string][]method{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || !fn.Name.IsExported() {
					continue
				}
				recv := receiverType(fn.Recv)
				if !strings.HasSuffix(recv, "Service") {
					continue
				}
				byService[recv] = append(byService[recv], method{
					Name:      fn.Name.Name,
					Signature: signature(fset, fn),
					Doc:       firstDocLine(fn.Doc),
				})
			}
		}
	}

	s := surface{SDK: "go", Source: filepath.ToSlash(dir)}
	for name, methods := range byService {
		sort.Slice(methods, func(i, j int) bool { return methods[i].Name < methods[j].Name })
		s.Services = append(s.Services, service{Name: name, Methods: methods})
	}
	sort.Slice(s.Services, func(i, j int) bool { return s.Services[i].Name < s.Services[j].Name })

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Printf("go: %d services\n", len(s.Services))
}

func receiverType(recv *ast.FieldList) string {
	if len(recv.List) == 0 {
		return ""
	}
	switch t := recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return id.Name
		}
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// signature renders the parameter and result lists as written in the source.
func signature(fset *token.FileSet, fn *ast.FuncDecl) string {
	var sb strings.Builder
	sb.WriteString(fn.Name.Name)
	sb.WriteString(render(fset, fn.Type.Params, true))
	if res := render(fset, fn.Type.Results, false); res != "" {
		sb.WriteString(" ")
		sb.WriteString(res)
	}
	return sb.String()
}

func render(fset *token.FileSet, fl *ast.FieldList, alwaysParens bool) string {
	if fl == nil || len(fl.List) == 0 {
		if alwaysParens {
			return "()"
		}
		return ""
	}
	parts := make([]string, 0, len(fl.List))
	for _, f := range fl.List {
		var buf strings.Builder
		_ = printer.Fprint(&buf, fset, f.Type)
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
		if len(names) == 0 {
			parts = append(parts, buf.String())
			continue
		}
		parts = append(parts, strings.Join(names, ", ")+" "+buf.String())
	}
	joined := strings.Join(parts, ", ")
	if alwaysParens || len(parts) > 1 {
		return "(" + joined + ")"
	}
	return joined
}

func firstDocLine(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	for _, line := range strings.Split(doc.Text(), "\n") {
		if t := strings.TrimSpace(line); t != "" {
			return t
		}
	}
	return ""
}
