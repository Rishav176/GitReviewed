package models

import "time"

// WebhookPayload represents the incoming webhook from GitHub
type WebhookPayload struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

// PullRequest contains PR details from GitHub
type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
	Head      GitRef    `json:"head"`
	Base      GitRef    `json:"base"`
}

// Repository contains repo information
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    User   `json:"owner"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

// GitRef represents a git reference (branch)
type GitRef struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

// DiffFile represents a single file change in a PR
type DiffFile struct {
	Filename  string
	Status    string // "added", "modified", "removed"
	Additions int
	Deletions int
	Changes   int
	Patch     string // The actual diff content
}

// ScanResult contains the results of security scanning
type ScanResult struct {
	Found      bool
	Issues     []SecurityIssue
	ScannedAt  time.Time
	TotalFiles int
}

// SecurityIssue represents a detected security problem
type SecurityIssue struct {
	Type        string // e.g., "AWS Access Key", "GitHub Token"
	FilePath    string
	LineNumber  int
	Severity    string // "CRITICAL", "HIGH", "MEDIUM", "LOW"
	Description string
	Pattern     string // Which pattern matched
}

// ReviewContext contains all info needed for a review
type ReviewContext struct {
	Repository  Repository
	PullRequest PullRequest
	DiffFiles   []DiffFile
	ScanResult  ScanResult
}

// SlackMessage represents the structure we'll send to Slack
type SlackMessage struct {
	Channel     string
	Text        string
	Blocks      interface{} // Slack Block Kit blocks
	ThreadTS    string      // For threading messages
	UnfurlLinks bool
}