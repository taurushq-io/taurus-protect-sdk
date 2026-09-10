package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The SDK must log metadata only — an id, a resource, a reason. Never a payload,
// token or key. That rule is stated on Logger and Field, but a doc comment does not
// stop the next person adding a log line, and this SDK's logs reach a consumer's
// disk under a policy that forbids payload logging outright.
//
// So the rule is checked mechanically: walk every s.logger call in this package and
// reject any argument that names a forbidden identifier. It fails when someone adds
// the log line, not in review.
func TestNoLoggerCallLeaksPayload(t *testing.T) {
	forbidden := []string{
		"PayloadAsString", "Payload",
		"Credentials", "ApiSecret", "APISecret", "Secret",
		"PrivateKey", "KeyPEM", "Token",
	}

	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	checked := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(".", e.Name()), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !isLoggerLevel(sel.Sel.Name) || !mentionsLogger(sel.X) {
				return true
			}
			checked++
			for _, arg := range call.Args {
				var buf strings.Builder
				renderExpr(&buf, arg)
				rendered := buf.String()
				for _, bad := range forbidden {
					if strings.Contains(rendered, bad) {
						t.Errorf("%s: logger call passes %q — SDK logs carry metadata only:\n  %s",
							fset.Position(call.Pos()), bad, rendered)
					}
				}
			}
			return true
		})
	}

	// A guard that inspects nothing passes for the wrong reason.
	if checked == 0 {
		t.Fatal("found no logger call sites to check; the guard is not looking where it should")
	}
}

func isLoggerLevel(name string) bool {
	switch name {
	case "Debug", "Info", "Warn", "Error":
		return true
	}
	return false
}

func mentionsLogger(x ast.Expr) bool {
	var buf strings.Builder
	renderExpr(&buf, x)
	return strings.Contains(strings.ToLower(buf.String()), "logger")
}

// renderExpr flattens an expression to its identifier names. Enough to spot a
// forbidden field being passed, without depending on a printer.
func renderExpr(buf *strings.Builder, e ast.Expr) {
	switch v := e.(type) {
	case *ast.Ident:
		buf.WriteString(v.Name)
	case *ast.SelectorExpr:
		renderExpr(buf, v.X)
		buf.WriteString("." + v.Sel.Name)
	case *ast.CallExpr:
		renderExpr(buf, v.Fun)
		for _, a := range v.Args {
			buf.WriteString(" ")
			renderExpr(buf, a)
		}
	case *ast.CompositeLit:
		for _, elt := range v.Elts {
			buf.WriteString(" ")
			renderExpr(buf, elt)
		}
	case *ast.KeyValueExpr:
		renderExpr(buf, v.Key)
		buf.WriteString(":")
		renderExpr(buf, v.Value)
	case *ast.BasicLit:
		buf.WriteString(v.Value)
	}
}
