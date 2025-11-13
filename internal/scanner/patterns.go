package scanner

import "regexp"

// SecretPattern defines a pattern for detecting secrets
type SecretPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Description string
	Severity    string
}

// GetDefaultPatterns returns the built-in secret detection patterns
func GetDefaultPatterns() []SecretPattern {
	return []SecretPattern{
		{
			Name:        "AWS Access Key ID",
			Pattern:     regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`),
			Description: "AWS Access Key ID detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "AWS Secret Access Key",
			Pattern:     regexp.MustCompile(`(?i)aws(.{0,20})?['\"][0-9a-zA-Z\/+]{40}['\"]`),
			Description: "AWS Secret Access Key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "GitHub Personal Access Token",
			Pattern:     regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
			Description: "GitHub personal access token detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "GitHub OAuth Token",
			Pattern:     regexp.MustCompile(`gho_[a-zA-Z0-9]{36}`),
			Description: "GitHub OAuth access token detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "GitHub App Token",
			Pattern:     regexp.MustCompile(`(ghu|ghs)_[a-zA-Z0-9]{36}`),
			Description: "GitHub App token detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "GitHub Refresh Token",
			Pattern:     regexp.MustCompile(`ghr_[a-zA-Z0-9]{36}`),
			Description: "GitHub refresh token detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "OpenAI API Key",
			Pattern:     regexp.MustCompile(`sk-[a-zA-Z0-9]{48}`),
			Description: "OpenAI API key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "Slack Token",
			Pattern:     regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`),
			Description: "Slack token detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "Slack Webhook",
			Pattern:     regexp.MustCompile(`https://hooks\.slack\.com/services/T[a-zA-Z0-9_]+/B[a-zA-Z0-9_]+/[a-zA-Z0-9_]+`),
			Description: "Slack webhook URL detected",
			Severity:    "HIGH",
		},
		{
			Name:        "Generic API Key",
			Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*['\"][a-zA-Z0-9]{20,}['\"]`),
			Description: "Generic API key pattern detected",
			Severity:    "HIGH",
		},
		{
			Name:        "Generic Secret",
			Pattern:     regexp.MustCompile(`(?i)(secret|password|passwd|pwd|token)\s*[:=]\s*['\"][^'\"]{8,}['\"]`),
			Description: "Generic secret pattern detected",
			Severity:    "MEDIUM",
		},
		{
			Name:        "Private Key",
			Pattern:     regexp.MustCompile(`-----BEGIN (RSA|DSA|EC|OPENSSH|PGP) PRIVATE KEY-----`),
			Description: "Private key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "Google API Key",
			Pattern:     regexp.MustCompile(`AIza[0-9A-Za-z\\-_]{35}`),
			Description: "Google API key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "Google OAuth",
			Pattern:     regexp.MustCompile(`[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com`),
			Description: "Google OAuth client ID detected",
			Severity:    "HIGH",
		},
		{
			Name:        "Stripe API Key",
			Pattern:     regexp.MustCompile(`(sk|pk)_(test|live)_[0-9a-zA-Z]{24,}`),
			Description: "Stripe API key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "Twilio API Key",
			Pattern:     regexp.MustCompile(`SK[0-9a-fA-F]{32}`),
			Description: "Twilio API key detected",
			Severity:    "CRITICAL",
		},
		{
			Name:        "JWT Token",
			Pattern:     regexp.MustCompile(`eyJ[A-Za-z0-9-_=]+\.eyJ[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`),
			Description: "JWT token detected",
			Severity:    "MEDIUM",
		},
		{
			Name:        "Database Connection String",
			Pattern:     regexp.MustCompile(`(?i)(mysql|postgres|mongodb|redis)://[^\s]+:[^\s]+@[^\s]+`),
			Description: "Database connection string with credentials detected",
			Severity:    "CRITICAL",
		},
	}
}

// ShouldIgnoreLine checks if a line should be ignored (e.g., comments, examples)
func ShouldIgnoreLine(line string) bool {
	ignorePatterns := []string{
		`(?i)example`,
		`(?i)sample`,
		`(?i)dummy`,
		`(?i)test`,
		`(?i)fake`,
		`(?i)placeholder`,
		`(?i)your[_-]?key[_-]?here`,
		`(?i)replace[_-]?with`,
		`(?i)TODO`,
		`(?i)FIXME`,
	}

	for _, pattern := range ignorePatterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			return true
		}
	}

	return false
}