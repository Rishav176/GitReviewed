package git

import (
	"context"

	"github.com/Rishav176/GitReviewed/internal/models"
)

// Client defines the interface for interacting with Git providers
type Client interface {
	// GetPRDiff fetches the diff for a pull request
	GetPRDiff(ctx context.Context, owner, repo string, prNumber int) ([]models.DiffFile, error)
	
	// VerifyWebhook verifies the webhook signature
	VerifyWebhook(payload []byte, signature string) bool
	
	// PostCommitStatus posts a status check to a commit
	PostCommitStatus(ctx context.Context, owner, repo, sha string, state, description, context string) error
}