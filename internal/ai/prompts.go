package ai

import (
	"fmt"
	"strings"

	"github.com/Rishav176/GitReviewed/internal/models"
)

const (
	MaxDiffLines     = 30   // REDUCED: Max lines per file
	MaxTotalFiles    = 3    // REDUCED: Max files to review
	MaxPromptLength  = 10000 // REDUCED: Max characters in prompt
)

// BuildReviewPrompt creates a prompt for code review with size limits
func BuildReviewPrompt(ctx models.ReviewContext) string {
	var prompt strings.Builder

	prompt.WriteString("You are an experienced code reviewer. Please review the following pull request changes.\n\n")
	prompt.WriteString(fmt.Sprintf("**Repository:** %s\n", ctx.Repository.FullName))
	prompt.WriteString(fmt.Sprintf("**PR Title:** %s\n", ctx.PullRequest.Title))
	prompt.WriteString(fmt.Sprintf("**Author:** %s\n\n", ctx.PullRequest.User.Login))

	prompt.WriteString("**Instructions:**\n")
	prompt.WriteString("1. Review the code for bugs, performance issues, and best practices\n")
	prompt.WriteString("2. Suggest improvements where applicable\n")
	prompt.WriteString("3. Point out any potential security issues\n")
	prompt.WriteString("4. If the code looks good, say so!\n")
	prompt.WriteString("5. Be constructive and helpful\n")
	prompt.WriteString("6. Keep your review concise\n\n")

	// Limit number of files
	filesToReview := ctx.DiffFiles
	if len(filesToReview) > MaxTotalFiles {
		prompt.WriteString(fmt.Sprintf("**Note:** Only reviewing first %d of %d files due to size limits.\n\n", MaxTotalFiles, len(filesToReview)))
		filesToReview = filesToReview[:MaxTotalFiles]
	}

	prompt.WriteString(fmt.Sprintf("**Changed Files:** (%d files)\n\n", len(filesToReview)))

	// Add each file's changes (with truncation)
	filesAdded := 0
	for _, file := range filesToReview {
		// Check if adding this file would exceed limit
		if prompt.Len() > MaxPromptLength {
			prompt.WriteString(fmt.Sprintf("\n**Note:** Stopped at %d files due to length limit.\n", filesAdded))
			break
		}

		prompt.WriteString(fmt.Sprintf("### File: `%s` (%s)\n", file.Filename, file.Status))
		prompt.WriteString(fmt.Sprintf("**Changes:** +%d additions, -%d deletions\n\n", file.Additions, file.Deletions))

		if file.Patch != "" {
			truncatedPatch := truncateDiff(file.Patch, MaxDiffLines)
			prompt.WriteString("```diff\n")
			prompt.WriteString(truncatedPatch)
			prompt.WriteString("\n```\n\n")
		}

		filesAdded++
	}

	prompt.WriteString("\n**Please provide your code review in 2-3 paragraphs:**")

	return prompt.String()
}

// truncateDiff truncates a diff to a maximum number of lines
func truncateDiff(diff string, maxLines int) string {
	lines := strings.Split(diff, "\n")
	
	if len(lines) <= maxLines {
		return diff
	}

	truncated := strings.Join(lines[:maxLines], "\n")
	truncated += fmt.Sprintf("\n... (truncated %d more lines)", len(lines)-maxLines)
	
	return truncated
}