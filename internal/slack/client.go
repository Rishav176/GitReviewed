package slack

import (
	"fmt"

	"github.com/Rishav176/GitReviewed/internal/models"
	"github.com/slack-go/slack"
)

// Client handles Slack API interactions
type Client struct {
	api            *slack.Client
	defaultChannel string
}

// NewClient creates a new Slack client
func NewClient(token, defaultChannel string) *Client {
	return &Client{
		api:            slack.New(token),
		defaultChannel: defaultChannel,
	}
}

// SendSecurityAlert sends a security alert about found secrets
func (c *Client) SendSecurityAlert(ctx models.ReviewContext) error {
	blocks := BuildSecurityAlertBlocks(ctx)

	_, _, err := c.api.PostMessage(
		c.defaultChannel,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionText("Security Alert: Secrets detected in PR", false),
	)

	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}

	return nil
}

// SendReviewComplete sends a message when review is complete with no issues
func (c *Client) SendReviewComplete(ctx models.ReviewContext) error {
	blocks := BuildReviewCompleteBlocks(ctx)

	_, _, err := c.api.PostMessage(
		c.defaultChannel,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionText("PR Review Complete: No issues found", false),
	)

	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}

	return nil
}

// TestConnection tests the Slack connection
func (c *Client) TestConnection() error {
	_, err := c.api.AuthTest()
	if err != nil {
		return fmt.Errorf("slack authentication failed: %w", err)
	}
	return nil
}