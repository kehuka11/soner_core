package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Service handles git operations.
type Service struct {
	baseDir string
}

// NewService creates a new Git service.
func NewService(baseDir string) *Service {
	return &Service{baseDir: baseDir}
}

// Clone clones a repository to a temporary directory and returns the path.
func (s *Service) Clone(repoURL string) (string, error) {
	// Create a unique directory name based on the repo URL or timestamp
	// For simplicity, we'll use a simple name for now, or a temp dir.
	// In a real app, we'd manage this better.
	repoName := filepath.Base(repoURL)
	targetDir := filepath.Join(s.baseDir, repoName)

	// Clean up if exists (for MVP)
	os.RemoveAll(targetDir)

	cmd := exec.Command("git", "clone", repoURL, targetDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to clone repo: %s, output: %s", err, output)
	}

	return targetDir, nil
}
