package analysis

import (
	"go/ast"
	"go/token"
)

// ComplexityScore represents the cognitive complexity of a function or file.
type ComplexityScore struct {
	Score   int
	Details []Detail
}

// Detail represents a specific point in the code that contributed to the complexity.
type Detail struct {
	Line    int
	Message string
	Cost    int
}

// Analyzer handles the complexity analysis.
type Analyzer struct {
	fset *token.FileSet
}

// NewAnalyzer creates a new Analyzer.
func NewAnalyzer(fset *token.FileSet) *Analyzer {
	return &Analyzer{fset: fset}
}

// AnalyzeFunction calculates the cognitive complexity of a function.
func (a *Analyzer) AnalyzeFunction(fn *ast.FuncDecl) ComplexityScore {
	v := &visitor{
		fset:    a.fset,
		details: []Detail{},
	}
	ast.Walk(v, fn.Body)
	return ComplexityScore{
		Score:   v.score,
		Details: v.details,
	}
}

type visitor struct {
	fset    *token.FileSet
	score   int
	nesting int
	details []Detail
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	cost := 0
	nestingCost := v.nesting

	switch n := node.(type) {
	case *ast.IfStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "if", cost)
		v.nesting++

		// Handle else
		if n.Else != nil {
			if _, ok := n.Else.(*ast.BlockStmt); ok {
				v.addDetail(n.Else.Pos(), "else", 1)
				v.score += 1
			}
		}
	case *ast.ForStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "for", cost)
		v.nesting++
	case *ast.RangeStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "range", cost)
		v.nesting++
	case *ast.SwitchStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "switch", cost)
		v.nesting++
	case *ast.TypeSwitchStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "type switch", cost)
		v.nesting++
	case *ast.SelectStmt:
		cost = 1 + nestingCost
		v.addDetail(n.Pos(), "select", cost)
		v.nesting++
	case *ast.CaseClause:
		// Case clauses don't increment nesting level for their body in Cognitive Complexity?
		// Actually, the switch itself adds nesting. The cases are just branches.
		// But if we have if inside case, it counts nesting from switch.
		// We need to handle nesting decrement carefully.
		// For now, let's assume standard visitor traversal.
		// Wait, ast.Walk visits children. We need to decrement nesting after visiting children.
		// But Visit returns a Visitor. If we return v, it continues.
		// We need to wrap the visitor to handle post-visit (decrement).
		return &nestingVisitor{v: v, originalNesting: v.nesting}
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			// Boolean operators add +1 but no nesting penalty
			// We only increment if the left child is NOT the same operator (to handle sequences)
			isSequence := false
			if left, ok := n.X.(*ast.BinaryExpr); ok {
				if left.Op == n.Op {
					isSequence = true
				}
			}

			if !isSequence {
				v.addDetail(n.Pos(), "boolean operator", 1)
				v.score += 1
			}
		}
	}

	v.score += cost
	return v
}

func (v *visitor) addDetail(pos token.Pos, msg string, cost int) {
	if cost > 0 {
		line := v.fset.Position(pos).Line
		v.details = append(v.details, Detail{
			Line:    line,
			Message: msg,
			Cost:    cost,
		})
	}
}

type nestingVisitor struct {
	v               *visitor
	originalNesting int
}

func (nv *nestingVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		// End of children, restore nesting
		nv.v.nesting = nv.originalNesting
		return nil
	}
	return nv.v.Visit(node)
}
