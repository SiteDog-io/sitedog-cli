package detectors

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

// maskSecrets masks sensitive information in debug output
// This function identifies common patterns for API keys, tokens, secrets and passwords in URLs
// and replaces them with masked versions showing only first 4 and last 4 characters
func maskSecrets(text string) string {
	// Define patterns for different types of secrets
	patterns := []string{
		`(sk-[a-zA-Z0-9]{48})`,           // OpenAI API keys (sk-...)
		`(pk_[a-zA-Z0-9]{50,})`,          // Stripe publishable keys
		`(sk_[a-zA-Z0-9]{50,})`,          // Stripe secret keys
		`(xoxb-[a-zA-Z0-9-]+)`,           // Slack bot tokens
		`(xoxp-[a-zA-Z0-9-]+)`,           // Slack user tokens
		`(ghp_[a-zA-Z0-9]{36})`,          // GitHub personal access tokens
		`(gho_[a-zA-Z0-9]{36})`,          // GitHub OAuth tokens
		`(AIza[a-zA-Z0-9_-]{35})`,        // Google API keys
		`([a-zA-Z0-9]{32,64})`,           // Generic long alphanumeric tokens
		`([A-Za-z0-9+/]{40,}={0,2})`,     // Base64 encoded secrets
		`([a-f0-9]{64})`,                 // 64-character hex tokens
		`([a-f0-9]{40})`,                 // 40-character hex tokens (like GitHub SHA)
	}

	result := text

	// First, handle passwords in URLs (like https://user:password@host)
	urlPasswordPattern := `(https?://[^:]+:)([^@]+)(@[^/\s]+)`
	urlRe := regexp.MustCompile(urlPasswordPattern)
	result = urlRe.ReplaceAllStringFunc(result, func(match string) string {
		parts := urlRe.FindStringSubmatch(match)
		if len(parts) == 4 {
			prefix := parts[1]    // https://user:
			password := parts[2]  // password
			suffix := parts[3]    // @host

			maskedPassword := password
			if len(password) > 8 {
				maskedPassword = password[:2] + "***" + password[len(password)-2:]
			}
			return prefix + maskedPassword + suffix
		}
		return match
	})

	// Then handle other secret patterns
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllStringFunc(result, func(match string) string {
			// Don't mask very short strings (likely not secrets)
			if len(match) <= 8 {
				return match
			}
			// Show first 4 and last 4 characters, mask the middle
			return match[:4] + "***" + match[len(match)-4:]
		})
	}

	return result
}

// adjustConfidenceForFileAge adjusts confidence based on file modification time
// Reduces confidence proportionally based on file age:
// - Files < 6 months: no reduction
// - Files 6-12 months: reduce by 5-15%
// - Files 1-2 years: reduce by 15-35%
// - Files 2-3 years: reduce by 35-50%
// - Files > 3 years: reduce by 50-70%
func adjustConfidenceForFileAge(baseConfidence float64, filePath string) float64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// If we can't get file info, return original confidence
		return baseConfidence
	}

	// Calculate age of file in months
	fileAge := time.Since(fileInfo.ModTime())
	ageInMonths := fileAge.Hours() / (24 * 30) // Approximate months

	var reductionFactor float64

	if ageInMonths < 6 {
		// Files less than 6 months old: no reduction
		reductionFactor = 1.0
	} else if ageInMonths < 12 {
		// 6-12 months: linear reduction from 0% to 15%
		progress := (ageInMonths - 6) / 6 // 0 to 1
		reductionFactor = 1.0 - (progress * 0.15)
	} else if ageInMonths < 24 {
		// 1-2 years: linear reduction from 15% to 35%
		progress := (ageInMonths - 12) / 12 // 0 to 1
		reductionFactor = 0.85 - (progress * 0.20) // 0.85 to 0.65
	} else if ageInMonths < 36 {
		// 2-3 years: linear reduction from 35% to 50%
		progress := (ageInMonths - 24) / 12 // 0 to 1
		reductionFactor = 0.65 - (progress * 0.15) // 0.65 to 0.50
	} else {
		// 3+ years: linear reduction from 50% to 70%, capped
		progress := (ageInMonths - 36) / 12 // 0 to 1 for year 4
		if progress > 1 {
			progress = 1 // Cap at maximum reduction
		}
		reductionFactor = 0.50 - (progress * 0.20) // 0.50 to 0.30
	}

	adjustedConfidence := baseConfidence * reductionFactor

	// Ensure confidence doesn't go below a reasonable minimum
	if adjustedConfidence < 0.25 {
		adjustedConfidence = 0.25
	}

	return adjustedConfidence
}

// formatFileAge formats file age for debug output
func formatFileAge(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "unknown age"
	}

	modTime := fileInfo.ModTime()
	age := time.Since(modTime)

	if age > 365*24*time.Hour {
		years := int(age.Hours() / (365 * 24))
		return modTime.Format("2006-01-02") + " (" + formatDuration(years, "year") + " old)"
	} else if age > 30*24*time.Hour {
		months := int(age.Hours() / (30 * 24))
		return modTime.Format("2006-01-02") + " (" + formatDuration(months, "month") + " old)"
	} else if age > 24*time.Hour {
		days := int(age.Hours() / 24)
		return modTime.Format("2006-01-02") + " (" + formatDuration(days, "day") + " old)"
	} else {
		return modTime.Format("2006-01-02") + " (recent)"
	}
}

// formatDuration formats duration with proper pluralization
func formatDuration(count int, unit string) string {
	if count == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", count, unit)
}