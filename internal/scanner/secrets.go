package scanner

import (
	"bufio"
	"strings"

	"github.com/Rishav176/GitReviewed/internal/models"
)

// Scanner handles secret detection
type Scanner struct {
	patterns []SecretPattern
}

// NewScanner creates a new scanner with default patterns
func NewScanner() *Scanner {
	return &Scanner{
		patterns: GetDefaultPatterns(),
	}
}

// NewScannerWithPatterns creates a scanner with custom patterns
func NewScannerWithPatterns(patterns []SecretPattern) *Scanner {
	return &Scanner{
		patterns: patterns,
	}
}

// ScanDiff scans a diff for secrets
func (s *Scanner) ScanDiff(diff string, filename string) []models.SecurityIssue {
	var issues []models.SecurityIssue

	scanner := bufio.NewScanner(strings.NewReader(diff))
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip lines that are removals (start with -)
		if strings.HasPrefix(line, "-") {
			continue
		}

		// Skip lines that should be ignored
		if ShouldIgnoreLine(line) {
			continue
		}

		// Check against all patterns
		for _, pattern := range s.patterns {
			if pattern.Pattern.MatchString(line) {
				issues = append(issues, models.SecurityIssue{
					Type:        pattern.Name,
					FilePath:    filename,
					LineNumber:  lineNumber,
					Severity:    pattern.Severity,
					Description: pattern.Description,
					Pattern:     pattern.Name,
				})
			}
		}
	}

	return issues
}

// ScanFiles scans multiple diff files
func (s *Scanner) ScanFiles(files []models.DiffFile) models.ScanResult {
	var allIssues []models.SecurityIssue

	for _, file := range files {
		issues := s.ScanDiff(file.Patch, file.Filename)
		allIssues = append(allIssues, issues...)
	}

	return models.ScanResult{
		Found:      len(allIssues) > 0,
		Issues:     allIssues,
		TotalFiles: len(files),
	}
}