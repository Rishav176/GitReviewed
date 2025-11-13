package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Rishav176/GitReviewed/internal/config"
	"github.com/Rishav176/GitReviewed/internal/git"
	"github.com/Rishav176/GitReviewed/internal/models"
	"github.com/Rishav176/GitReviewed/internal/scanner"
	"github.com/Rishav176/GitReviewed/internal/slack"
)

// WebhookHandler handles incoming GitHub webhooks
type WebhookHandler struct {
	config       *config.Config
	gitClient    git.Client
	slackClient  *slack.Client
	secretScanner *scanner.Scanner
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		config:        cfg,
		gitClient:     git.NewGitHubClient(cfg.GitHubToken, cfg.WebhookSecret),
		slackClient:   slack.NewClient(cfg.SlackToken, cfg.SlackChannel),
		secretScanner: scanner.NewScanner(),
	}
}

// HandleWebhook processes incoming GitHub webhook events
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify webhook signature
	signature := r.Header.Get("X-Hub-Signature-256")
	if !h.gitClient.VerifyWebhook(body, signature) {
		log.Printf("Invalid webhook signature")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the event type
	eventType := r.Header.Get("X-GitHub-Event")
	log.Printf("Received GitHub event: %s", eventType)

	// We only care about pull request events
	if eventType != "pull_request" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event ignored"))
		return
	}

	// Parse the webhook payload
	var payload models.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Only process opened or synchronize (new commits) actions
	if payload.Action != "opened" && payload.Action != "synchronize" {
		log.Printf("Ignoring action: %s", payload.Action)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Action ignored"))
		return
	}

	// Process the PR asynchronously
	go h.processPullRequest(payload)

	// Respond immediately
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

// processPullRequest handles the actual PR review
func (h *WebhookHandler) processPullRequest(payload models.WebhookPayload) {
	ctx := context.Background()

	log.Printf("Processing PR #%d from %s/%s",
		payload.PullRequest.Number,
		payload.Repository.Owner.Login,
		payload.Repository.Name,
	)

	// Fetch PR diff
	owner := payload.Repository.Owner.Login
	repo := payload.Repository.Name
	prNumber := payload.PullRequest.Number

	diffFiles, err := h.gitClient.GetPRDiff(ctx, owner, repo, prNumber)
	if err != nil {
		log.Printf("Error fetching PR diff: %v", err)
		return
	}

	log.Printf("Fetched %d files from PR #%d", len(diffFiles), prNumber)

	// Scan for secrets
	scanResult := h.secretScanner.ScanFiles(diffFiles)
	scanResult.ScannedAt = time.Now()

	log.Printf("Scan complete: found %d issues", len(scanResult.Issues))

	// Build review context
	reviewCtx := models.ReviewContext{
		Repository:  payload.Repository,
		PullRequest: payload.PullRequest,
		DiffFiles:   diffFiles,
		ScanResult:  scanResult,
	}

	// Send notification to Slack
	if scanResult.Found {
		log.Printf("Sending security alert to Slack")
		if err := h.slackClient.SendSecurityAlert(reviewCtx); err != nil {
			log.Printf("Error sending Slack alert: %v", err)
		}
	} else {
		log.Printf("Sending review complete message to Slack")
		if err := h.slackClient.SendReviewComplete(reviewCtx); err != nil {
			log.Printf("Error sending Slack message: %v", err)
		}
	}

	log.Printf("Completed processing PR #%d", prNumber)
}

// HealthCheck handles health check requests
func (h *WebhookHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TestSlack tests the Slack connection
func (h *WebhookHandler) TestSlack(w http.ResponseWriter, r *http.Request) {
	if err := h.slackClient.TestConnection(); err != nil {
		http.Error(w, fmt.Sprintf("Slack connection failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Slack connection successful"))
}