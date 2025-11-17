package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Rishav176/GitReviewed/internal/models"
	"google.golang.org/genai"
)

// Client handles AI API interactions using Google's official SDK
type Client struct {
	client *genai.Client
}

// NewClient creates a new AI client using the official Google SDK
func NewClient(apiKey string) *Client {
	ctx := context.Background()
	
	// Set API key as environment variable (SDK reads from GEMINI_API_KEY)
	os.Setenv("GEMINI_API_KEY", apiKey)
	
	// Create client (nil means it will use GEMINI_API_KEY from environment)
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		return nil
	}

	return &Client{
		client: client,
	}
}

// TestConnection tests the Gemini API with a simple request
func (c *Client) TestConnection() error {
	ctx := context.Background()

	result, err := c.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("Say hello in one word"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("API test failed: %w", err)
	}

	if result.Text() == "" {
		return fmt.Errorf("empty response from API")
	}

	return nil
}

// ReviewSingleFile reviews a single file
func (c *Client) ReviewSingleFile(filename string, patch string, additions, deletions int) (string, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf(`You are an experienced code reviewer. Review this single file change.

**File:** %s
**Changes:** +%d additions, -%d deletions

**Diff:**
`+"```diff\n%s\n```"+`

**Instructions:**
1. Review for bugs, performance issues, and best practices
2. Suggest specific improvements with line references if possible
3. Point out security issues
4. If the code looks good, briefly say so
5. Be concise - max 3-4 sentences per issue

**Your review:**`, filename, additions, deletions, patch)

	log.Printf("📊 Prompt size: %d characters", len(prompt))

	result, err := c.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	return result.Text(), nil
}

// ReviewCodeByFile reviews each file individually and combines results
func (c *Client) ReviewCodeByFile(ctx models.ReviewContext) (string, error) {
	var allReviews strings.Builder

	allReviews.WriteString(fmt.Sprintf("**PR Review for #%d: %s**\n\n", ctx.PullRequest.Number, ctx.PullRequest.Title))

	filesReviewed := 0
	filesFailed := 0

	// Review each file individually
	for i, file := range ctx.DiffFiles {
		// Skip binary files or files without patches
		if file.Patch == "" {
			continue
		}

		// Truncate very large diffs
		patch := file.Patch
		if len(patch) > 5000 { // ~100-150 lines of diff
			lines := strings.Split(patch, "\n")
			if len(lines) > 100 {
				patch = strings.Join(lines[:100], "\n")
				patch += fmt.Sprintf("\n... (truncated %d lines)", len(lines)-100)
			}
		}

		log.Printf("Reviewing file %d/%d: %s", i+1, len(ctx.DiffFiles), file.Filename)

		review, err := c.ReviewSingleFile(file.Filename, patch, file.Additions, file.Deletions)
		if err != nil {
			log.Printf("Failed to review %s: %v", file.Filename, err)
			allReviews.WriteString(fmt.Sprintf("\n### %s\n", file.Filename))
			allReviews.WriteString("_Could not review this file due to API error_\n\n")
			filesFailed++
			continue
		}

		allReviews.WriteString(fmt.Sprintf("\n### %s\n", file.Filename))
		allReviews.WriteString(review)
		allReviews.WriteString("\n\n")
		filesReviewed++

		// Rate limiting: wait 2 seconds between requests
		if i < len(ctx.DiffFiles)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	if filesReviewed == 0 {
		return "", fmt.Errorf("failed to review any files (%d failed)", filesFailed)
	}

	// Add overall summary
	allReviews.WriteString("\n---\n")
	allReviews.WriteString(fmt.Sprintf("**Summary:** Reviewed %d/%d file(s) successfully\n",
		filesReviewed,
		len(ctx.DiffFiles)))

	return allReviews.String(), nil
}