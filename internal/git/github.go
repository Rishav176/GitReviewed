package git

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Rishav176/GitReviewed/internal/models"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// GitHubClient implements the Client interface for GitHub
type GitHubClient struct {
	client        *github.Client
	webhookSecret string
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(token, webhookSecret string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{
		client:        github.NewClient(tc),
		webhookSecret: webhookSecret,
	}
}

// PostCommitStatus posts a status check to a commit
func (g *GitHubClient) PostCommitStatus(ctx context.Context, owner, repo, sha string, state, description, context string) error {
	status := &github.RepoStatus{
		State:       github.String(state),
		Description: github.String(description),
		Context:     github.String(context),
	}

	_, _, err := g.client.Repositories.CreateStatus(ctx, owner, repo, sha, status)
	if err != nil {
		return fmt.Errorf("failed to post commit status: %w", err)
	}

	return nil
}

// GetPRDiff fetches the diff for a pull request
func (g *GitHubClient) GetPRDiff(ctx context.Context, owner, repo string, prNumber int) ([]models.DiffFile, error) {
	// Get the list of files changed in the PR
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allFiles []models.DiffFile

	for {
		files, resp, err := g.client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch PR files: %w", err)
		}

		for _, file := range files {
			diffFile := models.DiffFile{
				Filename:  file.GetFilename(),
				Status:    file.GetStatus(),
				Additions: file.GetAdditions(),
				Deletions: file.GetDeletions(),
				Changes:   file.GetChanges(),
				Patch:     file.GetPatch(),
			}
			allFiles = append(allFiles, diffFile)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allFiles, nil
}

// VerifyWebhook verifies the GitHub webhook signature
func (g *GitHubClient) VerifyWebhook(payload []byte, signature string) bool {
	// GitHub sends the signature as "sha256=<signature>"
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	// Remove the "sha256=" prefix
	signature = strings.TrimPrefix(signature, "sha256=")

	// Compute HMAC
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GetPRInfo fetches basic PR information (useful for additional context)
func (g *GitHubClient) GetPRInfo(ctx context.Context, owner, repo string, prNumber int) (*models.PullRequest, error) {
	pr, _, err := g.client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR info: %w", err)
	}

	return &models.PullRequest{
		Number:    pr.GetNumber(),
		Title:     pr.GetTitle(),
		HTMLURL:   pr.GetHTMLURL(),
		State:     pr.GetState(),
		CreatedAt: pr.GetCreatedAt().Time,
		UpdatedAt: pr.GetUpdatedAt().Time,
		User: models.User{
			Login:     pr.GetUser().GetLogin(),
			ID:        pr.GetUser().GetID(),
			AvatarURL: pr.GetUser().GetAvatarURL(),
		},
		Head: models.GitRef{
			Ref: pr.GetHead().GetRef(),
			SHA: pr.GetHead().GetSHA(),
		},
		Base: models.GitRef{
			Ref: pr.GetBase().GetRef(),
			SHA: pr.GetBase().GetSHA(),
		},
	}, nil
}