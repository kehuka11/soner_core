package analysis

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAnalyzeFunction(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Simple function",
			code: `package main
				func foo() {
					println("hello")
				}`,
			expected: 0,
		},
		{
			name: "Single if",
			code: `package main
				func foo() {
					if true {
						println("true")
					}
				}`,
			expected: 1,
		},
		{
			name: "If else",
			code: `package main
				func foo() {
					if true {
						println("true")
					} else {
						println("false")
					}
				}`,
			expected: 2, // if (+1) + else (+1) = 2. Wait, else is +1? Whitepaper says "else if, else, default: +1".
		},
		{
			name: "Nested if",
			code: `package main
				func foo() {
					if true {
						if true {
							println("nested")
						}
					}
				}`,
			expected: 3, // if (+1) + nested if (+1 + 1 nesting) = 3
		},
		{
			name: "For loop",
			code: `package main
				func foo() {
					for i := 0; i < 10; i++ {
						println(i)
					}
				}`,
			expected: 1,
		},
		{
			name: "Nested for loop",
			code: `package main
				func foo() {
					for i := 0; i < 10; i++ {
						for j := 0; j < 10; j++ {
							println(j)
						}
					}
				}`,
			expected: 3, // for (+1) + nested for (+1 + 1 nesting) = 3
		},
		{
			name: "Switch case",
			code: `package main
				func foo() {
					switch 1 {
					case 1:
						println("one")
					case 2:
						println("two")
					default:
						println("default")
					}
				}`,
			expected: 1, // switch (+1). Cases don't increment.
		},
		{
			name: "Switch with nested if",
			code: `package main
				func foo() {
					switch 1 {
					case 1:
						if true {
							println("nested")
						}
					}
				}`,
			expected: 3, // switch (+1) + if (+1 + 1 nesting) = 3
		},
		{
			name: "Boolean operators",
			code: `package main
				func foo() {
					if true && true {
						println("true")
					}
				}`,
			expected: 2, // if (+1) + && (+1) = 2
		},
		{
			name: "Boolean sequence",
			code: `package main
				func foo() {
					if true && true && true {
						println("true")
					}
				}`,
			expected: 2, // if (+1) + && sequence (+1) = 2. Currently implementation might count each &&.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			var fn *ast.FuncDecl
			for _, decl := range node.Decls {
				if f, ok := decl.(*ast.FuncDecl); ok {
					fn = f
					break
				}
			}

			if fn == nil {
				t.Fatal("No function found in test code")
			}

			analyzer := NewAnalyzer(fset)
			score := analyzer.AnalyzeFunction(fn)

			if score.Score != tt.expected {
				t.Errorf("Expected score %d, got %d", tt.expected, score.Score)
				for _, d := range score.Details {
					t.Logf("Detail: %s (Line %d) +%d", d.Message, d.Line, d.Cost)
				}
			}
		})
	}
}
