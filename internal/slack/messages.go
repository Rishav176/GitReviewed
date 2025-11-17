package slack

import (
	"fmt"

	"github.com/Rishav176/GitReviewed/internal/models"
	"github.com/slack-go/slack"
)

// BuildSecurityAlertBlocks creates Slack blocks for security alerts
func BuildSecurityAlertBlocks(ctx models.ReviewContext) []slack.Block {
	blocks := []slack.Block{}

	// Header
	headerText := slack.NewTextBlockObject("mrkdwn", 
		fmt.Sprintf(":rotating_light: *Security Alert: Secrets Detected* :rotating_light:"), 
		false, false)
	headerBlock := slack.NewSectionBlock(headerText, nil, nil)
	blocks = append(blocks, headerBlock)

	// PR Information
	prInfoText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Repository:* %s\n*PR #%d:* <%s|%s>\n*Author:* %s",
			ctx.Repository.FullName,
			ctx.PullRequest.Number,
			ctx.PullRequest.HTMLURL,
			ctx.PullRequest.Title,
			ctx.PullRequest.User.Login,
		),
		false, false)
	prInfoBlock := slack.NewSectionBlock(prInfoText, nil, nil)
	blocks = append(blocks, prInfoBlock)

	// Divider
	blocks = append(blocks, slack.NewDividerBlock())

	// Issues summary
	summaryText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Found %d security issue(s) across %d file(s)*",
			len(ctx.ScanResult.Issues),
			ctx.ScanResult.TotalFiles,
		),
		false, false)
	summaryBlock := slack.NewSectionBlock(summaryText, nil, nil)
	blocks = append(blocks, summaryBlock)

	// Group issues by severity
	criticalIssues := []models.SecurityIssue{}
	highIssues := []models.SecurityIssue{}
	mediumIssues := []models.SecurityIssue{}
	lowIssues := []models.SecurityIssue{}

	for _, issue := range ctx.ScanResult.Issues {
		switch issue.Severity {
		case "CRITICAL":
			criticalIssues = append(criticalIssues, issue)
		case "HIGH":
			highIssues = append(highIssues, issue)
		case "MEDIUM":
			mediumIssues = append(mediumIssues, issue)
		case "LOW":
			lowIssues = append(lowIssues, issue)
		}
	}

	// Add issues by severity
	if len(criticalIssues) > 0 {
		blocks = append(blocks, buildIssueSection("CRITICAL", "🔴", criticalIssues)...)
	}
	if len(highIssues) > 0 {
		blocks = append(blocks, buildIssueSection("HIGH", "🟠", highIssues)...)
	}
	if len(mediumIssues) > 0 {
		blocks = append(blocks, buildIssueSection("MEDIUM", "🟡", mediumIssues)...)
	}
	if len(lowIssues) > 0 {
		blocks = append(blocks, buildIssueSection("LOW", "🟢", lowIssues)...)
	}

	// Divider
	blocks = append(blocks, slack.NewDividerBlock())

	// Action required
	actionText := slack.NewTextBlockObject("mrkdwn",
		"*⚠️ Action Required:* Please remove these secrets before merging!",
		false, false)
	actionBlock := slack.NewSectionBlock(actionText, nil, nil)
	blocks = append(blocks, actionBlock)

	// Button to view PR
	buttonText := slack.NewTextBlockObject("plain_text", "View Pull Request", false, false)
	button := slack.NewButtonBlockElement("view_pr", "view_pr", buttonText)
	button.URL = ctx.PullRequest.HTMLURL
	actionBlock2 := slack.NewActionBlock("pr_actions", button)
	blocks = append(blocks, actionBlock2)

	return blocks
}

// BuildAIReviewBlocks creates Slack blocks for AI code review
func BuildAIReviewBlocks(ctx models.ReviewContext, aiReview string) []slack.Block {
	blocks := []slack.Block{}

	// Header
	headerText := slack.NewTextBlockObject("mrkdwn",
		":robot_face: *AI Code Review*",
		false, false)
	headerBlock := slack.NewSectionBlock(headerText, nil, nil)
	blocks = append(blocks, headerBlock)

	// PR Information
	prInfoText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Repository:* %s\n*PR #%d:* <%s|%s>\n*Author:* %s",
			ctx.Repository.FullName,
			ctx.PullRequest.Number,
			ctx.PullRequest.HTMLURL,
			ctx.PullRequest.Title,
			ctx.PullRequest.User.Login,
		),
		false, false)
	prInfoBlock := slack.NewSectionBlock(prInfoText, nil, nil)
	blocks = append(blocks, prInfoBlock)

	// Divider
	blocks = append(blocks, slack.NewDividerBlock())

	// AI Review (split into chunks if too long)
	reviewText := slack.NewTextBlockObject("mrkdwn", aiReview, false, false)
	reviewBlock := slack.NewSectionBlock(reviewText, nil, nil)
	blocks = append(blocks, reviewBlock)

	// Divider
	blocks = append(blocks, slack.NewDividerBlock())

	// Button to view PR
	buttonText := slack.NewTextBlockObject("plain_text", "View Pull Request", false, false)
	button := slack.NewButtonBlockElement("view_pr", "view_pr", buttonText)
	button.URL = ctx.PullRequest.HTMLURL
	actionBlock := slack.NewActionBlock("pr_actions", button)
	blocks = append(blocks, actionBlock)

	return blocks
}

// buildIssueSection creates a section for a specific severity level
func buildIssueSection(severity, emoji string, issues []models.SecurityIssue) []slack.Block {
	blocks := []slack.Block{}

	// Severity header
	headerText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("%s *%s Severity* (%d issue(s))", emoji, severity, len(issues)),
		false, false)
	headerBlock := slack.NewSectionBlock(headerText, nil, nil)
	blocks = append(blocks, headerBlock)

	// List each issue
	for i, issue := range issues {
		if i >= 5 { // Limit to 5 issues per severity to avoid message bloat
			remainingText := slack.NewTextBlockObject("mrkdwn",
				fmt.Sprintf("_... and %d more %s severity issue(s)_", len(issues)-5, severity),
				false, false)
			remainingBlock := slack.NewSectionBlock(remainingText, nil, nil)
			blocks = append(blocks, remainingBlock)
			break
		}

		issueText := slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("• *%s*\n  `%s` (Line %d)\n  _%s_",
				issue.Type,
				issue.FilePath,
				issue.LineNumber,
				issue.Description,
			),
			false, false)
		issueBlock := slack.NewSectionBlock(issueText, nil, nil)
		blocks = append(blocks, issueBlock)
	}

	return blocks
}

// BuildReviewCompleteBlocks creates Slack blocks for successful review
func BuildReviewCompleteBlocks(ctx models.ReviewContext) []slack.Block {
	blocks := []slack.Block{}

	// Header
	headerText := slack.NewTextBlockObject("mrkdwn",
		":white_check_mark: *PR Review Complete*",
		false, false)
	headerBlock := slack.NewSectionBlock(headerText, nil, nil)
	blocks = append(blocks, headerBlock)

	// PR Information
	prInfoText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Repository:* %s\n*PR #%d:* <%s|%s>\n*Author:* %s",
			ctx.Repository.FullName,
			ctx.PullRequest.Number,
			ctx.PullRequest.HTMLURL,
			ctx.PullRequest.Title,
			ctx.PullRequest.User.Login,
		),
		false, false)
	prInfoBlock := slack.NewSectionBlock(prInfoText, nil, nil)
	blocks = append(blocks, prInfoBlock)

	// Divider
	blocks = append(blocks, slack.NewDividerBlock())

	// Success message
	successText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*No security issues found!*\nScanned %d file(s) - all clear! ✨",
			ctx.ScanResult.TotalFiles,
		),
		false, false)
	successBlock := slack.NewSectionBlock(successText, nil, nil)
	blocks = append(blocks, successBlock)

	// Button to view PR
	buttonText := slack.NewTextBlockObject("plain_text", "View Pull Request", false, false)
	button := slack.NewButtonBlockElement("view_pr", "view_pr", buttonText)
	button.URL = ctx.PullRequest.HTMLURL
	actionBlock := slack.NewActionBlock("pr_actions", button)
	blocks = append(blocks, actionBlock)

	return blocks
}