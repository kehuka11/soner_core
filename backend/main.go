package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"soner_core/backend/internal/analysis"
	"soner_core/backend/internal/git"
)

type AnalyzeRequest struct {
	RepoURL string `json:"repoUrl"`
}

type AnalyzeResponse struct {
	Files []FileResult `json:"files"`
}

type FileResult struct {
	Path       string           `json:"path"`
	Functions  []FunctionResult `json:"functions"`
	Complexity int              `json:"complexity"`
}

type FunctionResult struct {
	Name      string            `json:"name"`
	Score     int               `json:"score"`
	Details   []analysis.Detail `json:"details"`
	StartLine int               `json:"startLine"`
	EndLine   int               `json:"endLine"`
	Source    string            `json:"source"`
}

func main() {
	// Ensure temp dir exists
	tmpDir := "./tmp_repos"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatal(err)
	}

	gitService := git.NewService(tmpDir)

	http.HandleFunc("/api/analyze", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS for development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AnalyzeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Analyzing repo: %s", req.RepoURL)
		repoPath, err := gitService.Clone(req.RepoURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to clone: %v", err), http.StatusInternalServerError)
			return
		}
		// defer os.RemoveAll(repoPath) // Keep for inspection for now

		results, err := analyzeRepo(repoPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to analyze: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AnalyzeResponse{Files: results})
	})

	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func analyzeRepo(path string) ([]FileResult, error) {
	var results []FileResult
	fset := token.NewFileSet()

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(filePath, ".go") {
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read %s: %v", filePath, err)
				return nil
			}

			node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
			if err != nil {
				log.Printf("Failed to parse %s: %v", filePath, err)
				return nil // Continue
			}

			analyzer := analysis.NewAnalyzer(fset)
			var fileFuncs []FunctionResult
			totalComplexity := 0

			for _, decl := range node.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					score := analyzer.AnalyzeFunction(fn)

					start := fset.Position(fn.Pos())
					end := fset.Position(fn.End())

					// Extract source code
					// Note: Position.Offset is 0-indexed byte offset
					source := string(content[start.Offset:end.Offset])

					fileFuncs = append(fileFuncs, FunctionResult{
						Name:      fn.Name.Name,
						Score:     score.Score,
						Details:   score.Details,
						StartLine: start.Line,
						EndLine:   end.Line,
						Source:    source,
					})
					totalComplexity += score.Score
				}
			}

			relPath, _ := filepath.Rel(path, filePath)
			results = append(results, FileResult{
				Path:       relPath,
				Functions:  fileFuncs,
				Complexity: totalComplexity,
			})
		}
		return nil
	})

	return results, err
}
